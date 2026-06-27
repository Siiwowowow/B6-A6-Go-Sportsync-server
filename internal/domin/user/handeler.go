// internal/user/handler.go
package user

import (
	"errors"
	"gotickets/internal/domin/user/dto"
	"net/http"

	"github.com/labstack/echo/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// CreateUser handles registration POST /api/v1/auth/register
func (h *Handler) CreateUser(c *echo.Context) error {
	var req dto.CreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "Invalid request payload",
			"errors":  err.Error(),
		})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "Validation failed",
			"errors":  err.Error(),
		})
	}
	response, err := h.service.CreateUser(req)
	if err != nil {
		if errors.Is(err, ErrorAlreadyExists) {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"success": false,
				"message": "User already exists",
				"errors":  err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to create user",
			"errors":  err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]any{
		"success": true,
		"message": "User registered successfully",
		"data":    response,
	})
}

// LoginUser handles login POST /api/v1/auth/login
func (h *Handler) LoginUser(c *echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "Invalid request payload",
			"errors":  err.Error(),
		})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "Validation failed",
			"errors":  err.Error(),
		})
	}

	response, err := h.service.LoginUser(req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			return c.JSON(http.StatusUnauthorized, map[string]any{
				"success": false,
				"message": "Invalid email or password",
				"errors":  err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to login",
			"errors":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "Login successful",
		"data":    response,
	})
}

// Getme handles GET /api/v1/auth/me (Protected)
func (h *Handler) Getme(c *echo.Context) error {
	userID := c.Get("user_id")
	userEmail := c.Get("user_email")
	userName := c.Get("user_name")
	userRole := c.Get("user_role")

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "User retrieved successfully",
		"data": map[string]any{
			"user_id": userID,
			"email":   userEmail,
			"name":    userName,
			"role":    userRole,
		},
	})
}
