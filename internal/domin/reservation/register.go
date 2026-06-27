// internal/domin/reservation/register.go
package reservation

import (
	"gotickets/internal/auth"
	"gotickets/internal/middleware"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, jwtService auth.JWTService) {
	repo := NewRepository(db)
	svc := NewService(repo)
	handler := NewHandler(svc)

	// Authenticated routes
	authGroup := e.Group("/api/v1/reservations")
	authGroup.Use(middleware.AuthMiddleware(jwtService))

	authGroup.POST("", handler.CreateReservation)
	authGroup.GET("/my-reservations", handler.GetMyReservations)
	authGroup.DELETE("/:id", handler.CancelReservation)

	// Admin-only route
	adminGroup := e.Group("/api/v1/reservations")
	adminGroup.Use(middleware.AuthMiddleware(jwtService))
	adminGroup.Use(middleware.RoleMiddleware("admin"))
	adminGroup.GET("", handler.GetAllReservations)
}
