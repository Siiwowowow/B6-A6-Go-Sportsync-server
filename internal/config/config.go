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
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
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
