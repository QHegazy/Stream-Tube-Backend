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
	Salt            string

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
			RedisPassword:   getRequiredEnv("REDIS_PASSWORD"),
			RedisPort:       getRequiredEnv("REDIS_PORT"),
			PostgresUser:    getRequiredEnv("POSTGRES_USER"),
			PostgresPassword:getRequiredEnv("POSTGRES_PASSWORD"),
			PostgresDB:      getRequiredEnv("POSTGRES_DB"),
			PostgresHost:    getRequiredEnv("POSTGRES_HOST"),
			PostgresPort:    getRequiredEnv("POSTGRES_PORT"),
			SSLMode:         getRequiredEnv("SSL_MODE"),
			JWTSecretKey:    getRequiredEnv("SECRET_KEY"),
			ClientSide :     getRequiredEnv("CLIENT_SIDE"),
			Salt: 			 getRequiredEnv("SALT"),
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
