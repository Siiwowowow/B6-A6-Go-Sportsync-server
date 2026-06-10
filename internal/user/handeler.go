// internal/user/handler.go
package user

import (
	"errors"
	"gotickets/internal/httpResponse"
	"gotickets/internal/user/dto"
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
func (h *Handler) CreateUser(c *echo.Context) error {
	var req dto.CreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpResponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
			Details: err.Error()})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpResponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Validation failed",
			Details: err.Error()})
	}
	response, err := h.service.CreateUser(req)
	if err != nil {
		if errors.Is(err, ErrorAlreadyExists) {
			return c.JSON(http.StatusConflict, httpResponse.Error{
				Code:    http.StatusConflict,
				Message: "User already exists",
				Details: err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, httpResponse.Error{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create user",
			Details: err.Error(),
		})
	}
	return c.JSON(http.StatusCreated, response)

}
func (h *Handler) LoginUser(c *echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpResponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
			Details: err.Error()})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpResponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Validation failed",
			Details: err.Error()})
	}
	h.service.LoginUser(req)
	response, err := h.service.LoginUser(req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			return c.JSON(http.StatusUnauthorized, httpResponse.Error{
				Code:    http.StatusUnauthorized,
				Message: "Invalid email or password",
				Details: err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, httpResponse.Error{
			Code:    http.StatusInternalServerError,
			Message: "Failed to login",
			Details: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":    http.StatusOK,
		"message": "Login successful",
		"data":    response,
	})
}

func (h *Handler) Getme(c *echo.Context) error {
	userID := c.Get("user_id")
	userEmail := c.Get("user_email")
	userName := c.Get("user_name")

	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":    http.StatusOK,
		"message": "User retrieved successfully",
		"data": map[string]interface{}{
			"user_id": userID,
			"email":   userEmail,
			"name":    userName,
		},
	})
}
