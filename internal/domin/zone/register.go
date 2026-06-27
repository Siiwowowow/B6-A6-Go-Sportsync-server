// internal/domin/zone/register.go
package zone

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

	// Public routes
	e.GET("/api/v1/zones", handler.GetAllZones)
	e.GET("/api/v1/zones/:id", handler.GetZoneByID)

	// Admin-only routes
	adminGroup := e.Group("/api/v1/zones")
	adminGroup.Use(middleware.AuthMiddleware(jwtService))
	adminGroup.Use(middleware.RoleMiddleware("admin"))
	adminGroup.POST("", handler.CreateZone)
	adminGroup.PUT("/:id", handler.UpdateZone)
	adminGroup.DELETE("/:id", handler.DeleteZone)
}
