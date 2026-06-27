package middleware

import (
	"gotickets/internal/auth"
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
)

func AuthMiddleware(jwtService auth.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]any{
					"success": false,
					"message": "Missing Authorization header",
					"errors":  "Unauthorized access",
				})
			}
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]any{
					"success": false,
					"message": "Invalid Authorization header format",
					"errors":  "Unauthorized access",
				})
			}
			tokenString := parts[1]

			claims, err := jwtService.ValidateToken(tokenString)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]any{
					"success": false,
					"message": "Invalid token",
					"errors":  err.Error(),
				})
			}
			c.Set("user_id", claims.UserID)
			c.Set("user_email", claims.Email)
			c.Set("user_name", claims.Name)
			c.Set("user_role", claims.Role)
			return next(c)
		}
	}
}

// RoleMiddleware checks if the user has one of the required roles
func RoleMiddleware(requiredRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			userRole, ok := c.Get("user_role").(string)
			if !ok {
				return c.JSON(http.StatusForbidden, map[string]any{
					"success": false,
					"message": "Access forbidden: role not found",
					"errors":  "Forbidden access",
				})
			}
			for _, role := range requiredRoles {
				if userRole == role {
					return next(c)
				}
			}
			return c.JSON(http.StatusForbidden, map[string]any{
				"success": false,
				"message": "Access forbidden: insufficient permissions",
				"errors":  "Forbidden access",
			})
		}
	}
}
