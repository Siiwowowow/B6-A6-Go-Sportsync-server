package booking

import (
	"gotickets/internal/auth"
	"gotickets/internal/config"
	"gotickets/internal/event"
	"gotickets/internal/middleware"
	"time"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, cfg *config.Config) {
	bookingRepo := NewRepository(db)
	eventRepo := event.NewRepository(db)

	svc := NewService(bookingRepo, eventRepo)
	handler := NewHandler(svc)

	jwtService := auth.NewJWTService(cfg.JwtSecret, 24*time.Hour)

	api := e.Group("/api/v1/bookings", middleware.AuthMiddleware(jwtService))

	api.POST("", handler.CreateBooking)
	api.GET("/me", handler.GetMyBookings)

}
