package config

import (
	"os"
)

type Config struct {
	RedisPassword   string
	RedisPort      string
	PostgresUser   string
	PostgresPassword string
	PostgresDB     string
	PostgresHost   string
	PostgresPort   string
	SSLMode        string
}

func LoadConfig() *Config {
	return &Config{
		RedisPassword:   os.Getenv("REDIS_PASSWORD"),
		RedisPort:      os.Getenv("REDIS_PORT"),
		PostgresUser:   os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresDB:     os.Getenv("POSTGRES_DB"),
		PostgresHost:   os.Getenv("POSTGRES_HOST"),
		PostgresPort:   os.Getenv("POSTGRES_PORT"),
		SSLMode:        os.Getenv("SSL_MODE"),
	}
}
