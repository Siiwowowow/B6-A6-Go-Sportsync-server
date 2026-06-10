// internal/user/register.go
package user

import (
	"gotickets/internal/auth"
	"gotickets/internal/config"
	"gotickets/internal/middleware"
	"time"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, cfg *config.Config) {
	userRepository := NewRepository(db)
	jwtService := auth.NewJWTService(cfg.JwtSecret, 24*time.Hour)
	userService := NewService(userRepository, jwtService)
	userHandler := NewHandler(userService)
	api := e.Group("/api/v1/auth")
	api.POST("/register", userHandler.CreateUser)
	api.POST("/login", userHandler.LoginUser)

	// Protected route with authentication middleware
	protected := e.Group("/api/v1/auth")
	protected.Use(middleware.AuthMiddleware(jwtService))
	protected.GET("/me", userHandler.Getme)
}
