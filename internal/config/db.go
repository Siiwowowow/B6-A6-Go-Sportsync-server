// internal/config/db.go
package config

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectDB establishes a connection to the PostgreSQL database using GORM.
// It also configures a custom logger with a 3-second slow query threshold.
func ConnectDB(cfg *Config) *gorm.DB {
	dsn := cfg.Dsn

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             3 * time.Second,
			LogLevel:                  logger.Error,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
		Logger:         newLogger,
	})

	if err != nil {
		panic("failed to connect database")
	} else {
		println("Database connection successful")
	}
	return db
}