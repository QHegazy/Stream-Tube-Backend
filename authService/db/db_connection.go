package db

import (
	"authService/config"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Connections holds different connection pools for different operations
type Connections struct {
    ReadPool   *pgxpool.Pool
    CreatePool *pgxpool.Pool
    UpdatePool *pgxpool.Pool
    DeletePool *pgxpool.Pool
}

var pools *Connections

// InitDB initializes separate connection pools for different operations
func InitDB(cfg *config.Config) error {
    pools = &Connections{}
    
    // Create the base connection string
    connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
        cfg.PostgresUser,
        cfg.PostgresPassword,
        cfg.PostgresHost,
        cfg.PostgresPort,
        cfg.PostgresDB,
        cfg.SSLMode)

    // Initialize each pool with specific configurations
    var err error

    // Read Pool (optimized for reads)
    pools.ReadPool, err = createPool(connString, PoolConfig{
        MaxConns:         30,
        MinConns:         10,
        MaxConnIdleTime:  5 * time.Minute,
        MaxConnLifetime:  30 * time.Minute,
        HealthCheckPeriod: 30 * time.Second,
    })
    if err != nil {
        return fmt.Errorf("failed to create read pool: %v", err)
    }

    // Create Pool
    pools.CreatePool, err = createPool(connString, PoolConfig{
        MaxConns:         15,
        MinConns:         5,
        MaxConnIdleTime:  5 * time.Minute,
        MaxConnLifetime:  30 * time.Minute,
        HealthCheckPeriod: 30 * time.Second,
    })
    if err != nil {
        return fmt.Errorf("failed to create create pool: %v", err)
    }

    // Update Pool
    pools.UpdatePool, err = createPool(connString, PoolConfig{
        MaxConns:         15,
        MinConns:         5,
        MaxConnIdleTime:  5 * time.Minute,
        MaxConnLifetime:  30 * time.Minute,
        HealthCheckPeriod: 30 * time.Second,
    })
    if err != nil {
        return fmt.Errorf("failed to create update pool: %v", err)
    }

    // Delete Pool
    pools.DeletePool, err = createPool(connString, PoolConfig{
        MaxConns:         10,
        MinConns:         3,
        MaxConnIdleTime:  5 * time.Minute,
        MaxConnLifetime:  30 * time.Minute,
        HealthCheckPeriod: 30 * time.Second,
    })
    if err != nil {
        return fmt.Errorf("failed to create delete pool: %v", err)
    }

    log.Println("All database connection pools created successfully")
    return nil
}

// PoolConfig holds configuration for a single pool
type PoolConfig struct {
    MaxConns          int32
    MinConns          int32
    MaxConnIdleTime   time.Duration
    MaxConnLifetime   time.Duration
    HealthCheckPeriod time.Duration
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

// GetPools returns all connection pools
func GetPools() *Connections {
    return pools
}

// CloseAllPools closes all database connection pools
func CloseAllPools() {
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
    return map[string]*pgxpool.Stat{
        "read":   pools.ReadPool.Stat(),
        "create": pools.CreatePool.Stat(),
        "update": pools.UpdatePool.Stat(),
        "delete": pools.DeletePool.Stat(),
    }
}

// HealthCheck performs a health check on all pools
func HealthCheck(ctx context.Context) map[string]error {
    return map[string]error{
        "read":   pools.ReadPool.Ping(ctx),
        "create": pools.CreatePool.Ping(ctx),
        "update": pools.UpdatePool.Ping(ctx),
        "delete": pools.DeletePool.Ping(ctx),
    }
}