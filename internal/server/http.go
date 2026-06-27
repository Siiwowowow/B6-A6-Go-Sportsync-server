// internal/server/http.go
package server

import (
	"fmt"
	"gotickets/internal/auth"
	"gotickets/internal/config"
	"gotickets/internal/domin/reservation"
	"gotickets/internal/domin/user"
	"gotickets/internal/domin/zone"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"gorm.io/gorm"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.ErrBadRequest.Wrap(err)
	}
	return nil
}

func Start(db *gorm.DB, cfg *config.Config) {
	// Auto migrate database tables
	db.AutoMigrate(&user.User{})
	db.AutoMigrate(&zone.ParkingZone{})
	db.AutoMigrate(&reservation.Reservation{})

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.RequestLogger())

	e.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Hello,Go World!")
	})

	// Instantiate shared JWT service for route validation
	jwtService := auth.NewJWTService(cfg.JwtSecret, 24*time.Hour)

	user.RegisterRoutes(e, db, cfg)
	zone.RegisterRoutes(e, db, jwtService)
	reservation.RegisterRoutes(e, db, jwtService)

	port := fmt.Sprintf(":%s", config.LoadEnv().Port)
	serverURL := fmt.Sprintf("http://localhost:%s", config.LoadEnv().Port)
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("🚀 Server is running!")
	fmt.Println("📍 URL: " + serverURL)
	fmt.Println(strings.Repeat("=", 50) + "\n")
	if err := e.Start(port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
