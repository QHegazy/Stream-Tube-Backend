package db

import (
	"authService/config"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Connections holds different connection pools for different operations
type Connections struct {
	ReadPool   *pgxpool.Pool
	CreatePool *pgxpool.Pool
	UpdatePool *pgxpool.Pool
	DeletePool *pgxpool.Pool
	mu         sync.RWMutex
}

var (
	pools    *Connections
	initOnce sync.Once
)

// PoolConfig holds configuration for a single pool
type PoolConfig struct {
	MaxConns          int32
	MinConns          int32
	MaxConnIdleTime   time.Duration
	MaxConnLifetime   time.Duration
	HealthCheckPeriod time.Duration
}

// Default pool configurations
var (
	readPoolConfig = PoolConfig{
		MaxConns:          30,
		MinConns:         10,
		MaxConnIdleTime:   5 * time.Minute,
		MaxConnLifetime:   30 * time.Minute,
		HealthCheckPeriod: 30 * time.Second,
	}

	writePoolConfig = PoolConfig{
		MaxConns:          15,
		MinConns:         5,
		MaxConnIdleTime:   5 * time.Minute,
		MaxConnLifetime:   30 * time.Minute,
		HealthCheckPeriod: 30 * time.Second,
	}

	deletePoolConfig = PoolConfig{
		MaxConns:          10,
		MinConns:         3,
		MaxConnIdleTime:   5 * time.Minute,
		MaxConnLifetime:   30 * time.Minute,
		HealthCheckPeriod: 30 * time.Second,
	}
)

// InitDB initializes separate connection pools for different operations
func InitDB(cfg *config.Config) error {
	var initErr error
	
	initOnce.Do(func() {
		pools = &Connections{}
		
		// Create the base connection string
		connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
			cfg.PostgresUser,
			cfg.PostgresPassword,
			cfg.PostgresHost,
			cfg.PostgresPort,
			cfg.PostgresDB,
			cfg.SSLMode)

		// Initialize pools concurrently
		var wg sync.WaitGroup
		errChan := make(chan error, 4)

		// Initialize Read Pool
		wg.Add(1)
		go func() {
			defer wg.Done()
			var err error
			pools.ReadPool, err = createPool(connString, readPoolConfig)
			if err != nil {
				errChan <- fmt.Errorf("failed to create read pool: %v", err)
			}
		}()

		// Initialize Create Pool
		wg.Add(1)
		go func() {
			defer wg.Done()
			var err error
			pools.CreatePool, err = createPool(connString, writePoolConfig)
			if err != nil {
				errChan <- fmt.Errorf("failed to create create pool: %v", err)
			}
		}()

		// Initialize Update Pool
		wg.Add(1)
		go func() {
			defer wg.Done()
			var err error
			pools.UpdatePool, err = createPool(connString, writePoolConfig)
			if err != nil {
				errChan <- fmt.Errorf("failed to create update pool: %v", err)
			}
		}()

		// Initialize Delete Pool
		wg.Add(1)
		go func() {
			defer wg.Done()
			var err error
			pools.DeletePool, err = createPool(connString, deletePoolConfig)
			if err != nil {
				errChan <- fmt.Errorf("failed to create delete pool: %v", err)
			}
		}()

		// Wait for all pools to initialize
		wg.Wait()
		close(errChan)

		// Check for any errors
		for err := range errChan {
			if err != nil {
				initErr = err
				return
			}
		}

		log.Println("All database connection pools created successfully")
	})

	return initErr
}

// createPool creates a single connection pool with the given configuration
func createPool(connString string, cfg PoolConfig) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %v", err)
	}

	// Apply configuration
	poolConfig.MaxConns = cfg.MaxConns
	poolConfig.MinConns = cfg.MinConns
	poolConfig.MaxConnIdleTime = cfg.MaxConnIdleTime
	poolConfig.MaxConnLifetime = cfg.MaxConnLifetime
	poolConfig.HealthCheckPeriod = cfg.HealthCheckPeriod

	// Add connection lifecycle hooks for monitoring
	poolConfig.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
		// Add any connection validation logic here
		return true
	}

	poolConfig.AfterRelease = func(conn *pgx.Conn) bool {
		// Add any cleanup logic here
		return true
	}

	// Create the pool
	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}

	// Test the connection
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("unable to ping database: %v", err)
	}

	return pool, nil
}

// GetPools returns all connection pools with thread-safe access
func GetPools() *Connections {
	return pools
}

// CloseAllPools closes all database connection pools safely
func CloseAllPools() {
	if pools == nil {
		return
	}

	pools.mu.Lock()
	defer pools.mu.Unlock()

	if pools.ReadPool != nil {
		pools.ReadPool.Close()
	}
	if pools.CreatePool != nil {
		pools.CreatePool.Close()
	}
	if pools.UpdatePool != nil {
		pools.UpdatePool.Close()
	}
	if pools.DeletePool != nil {
		pools.DeletePool.Close()
	}
}

// GetStats returns statistics for all pools
func GetStats() map[string]*pgxpool.Stat {
	if pools == nil {
		return nil
	}

	pools.mu.RLock()
	defer pools.mu.RUnlock()

	return map[string]*pgxpool.Stat{
		"read":   pools.ReadPool.Stat(),
		"create": pools.CreatePool.Stat(),
		"update": pools.UpdatePool.Stat(),
		"delete": pools.DeletePool.Stat(),
	}
}

// HealthCheck performs a health check on all pools
func HealthCheck(ctx context.Context) map[string]error {
	if pools == nil {
		return map[string]error{"init": fmt.Errorf("pools not initialized")}
	}

	pools.mu.RLock()
	defer pools.mu.RUnlock()

	return map[string]error{
		"read":   pools.ReadPool.Ping(ctx),
		"create": pools.CreatePool.Ping(ctx),
		"update": pools.UpdatePool.Ping(ctx),
		"delete": pools.DeletePool.Ping(ctx),
	}
}

// RefreshPools attempts to recreate any unhealthy pools
func RefreshPools(ctx context.Context, cfg *config.Config) error {
	if pools == nil {
		return fmt.Errorf("pools not initialized")
	}

	healthStatus := HealthCheck(ctx)
	for poolName, err := range healthStatus {
		if err != nil {
			log.Printf("Unhealthy pool detected: %s, attempting to refresh", poolName)
			if refreshErr := refreshPool(ctx, cfg, poolName); refreshErr != nil {
				return fmt.Errorf("failed to refresh %s pool: %v", poolName, refreshErr)
			}
		}
	}
	return nil
}

// refreshPool recreates a specific pool
func refreshPool(ctx context.Context, cfg *config.Config, poolName string) error {
	pools.mu.Lock()
	defer pools.mu.Unlock()

	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDB,
		cfg.SSLMode)

	var newPool *pgxpool.Pool
	var err error

	switch poolName {
	case "read":
		newPool, err = createPool(connString, readPoolConfig)
	case "create", "update":
		newPool, err = createPool(connString, writePoolConfig)
	case "delete":
		newPool, err = createPool(connString, deletePoolConfig)
	default:
		return fmt.Errorf("unknown pool type: %s", poolName)
	}

	if err != nil {
		return err
	}

	// Replace the old pool with the new one
	switch poolName {
	case "read":
		if pools.ReadPool != nil {
			pools.ReadPool.Close()
		}
		pools.ReadPool = newPool
	case "create":
		if pools.CreatePool != nil {
			pools.CreatePool.Close()
		}
		pools.CreatePool = newPool
	case "update":
		if pools.UpdatePool != nil {
			pools.UpdatePool.Close()
		}
		pools.UpdatePool = newPool
	case "delete":
		if pools.DeletePool != nil {
			pools.DeletePool.Close()
		}
		pools.DeletePool = newPool
	}

	return nil
}