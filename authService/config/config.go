package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var (
	config     *Config
	configOnce sync.Once
)

type Config struct {
	// Database configs
	RedisPassword    string
	RedisPort       string
	PostgresUser    string
	PostgresPassword string
	PostgresDB      string
	PostgresHost    string
	PostgresPort    string
	SSLMode         string
	JWTSecretKey    string
	ClientSide      string

}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func GetConfig() *Config {
	configOnce.Do(func() {
		config = &Config{
			RedisPassword:   os.Getenv("REDIS_PASSWORD"),
			RedisPort:       os.Getenv("REDIS_PORT"),
			PostgresUser:    os.Getenv("POSTGRES_USER"),
			PostgresPassword:os.Getenv("POSTGRES_PASSWORD"),
			PostgresDB:      os.Getenv("POSTGRES_DB"),
			PostgresHost:    os.Getenv("POSTGRES_HOST"),
			PostgresPort:    os.Getenv("POSTGRES_PORT"),
			SSLMode:         os.Getenv("SSL_MODE"),
			JWTSecretKey:    getRequiredEnv("SECRET_KEY"),
			ClientSide :     getRequiredEnv("CLIENT_SIDE"),
		}
	})
	return config
}

func getRequiredEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(key + " environment variable is not set")
	}
	return value
}

 func GetClientSide() string {
	return GetConfig().ClientSide
}
