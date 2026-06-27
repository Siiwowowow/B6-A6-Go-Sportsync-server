// internal/config/config.go
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds system environment variables
type Config struct {
	Port      string
	Dsn       string
	JwtSecret string
}

// LoadEnv reads properties from .env and maps them to Config
func LoadEnv() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DSN")
	if dsn == "" {
		dsn = os.Getenv("DATABASE_URL")
	}

	return &Config{
		Port:      os.Getenv("PORT"),
		Dsn:       dsn,
		JwtSecret: os.Getenv("JWT_SECRET"),
	}
}
