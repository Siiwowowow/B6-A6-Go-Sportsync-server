// internal/domin/reservation/handler.go
package reservation

import (
	"errors"
	"gotickets/internal/domin/reservation/dto"
	"gotickets/internal/domin/zone"
	"net/http"

	"github.com/google/uuid"
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

// CreateReservation handles POST /api/v1/reservations (Authenticated)
func (h *Handler) CreateReservation(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{
			"success": false,
			"message": "User not authenticated",
			"errors":  "Unauthorized",
		})
	}

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

	response, err := h.service.CreateReservation(userID, req)
	if err != nil {
		if errors.Is(err, ErrZoneFull) {
			return c.JSON(http.StatusConflict, map[string]any{
				"success": false,
				"message": "Reservation failed",
				"errors":  err.Error(),
			})
		}
		if errors.Is(err, zone.ErrZoneNotFound) {
			return c.JSON(http.StatusNotFound, map[string]any{
				"success": false,
				"message": "Parking zone not found",
				"errors":  err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to confirm reservation",
			"errors":  err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]any{
		"success": true,
		"message": "Reservation confirmed successfully",
		"data":    response,
	})
}

// GetMyReservations handles GET /api/v1/reservations/my-reservations (Authenticated)
func (h *Handler) GetMyReservations(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{
			"success": false,
			"message": "User not authenticated",
			"errors":  "Unauthorized",
		})
	}

	response, err := h.service.GetMyReservations(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to retrieve reservations",
			"errors":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "My reservations retrieved successfully",
		"data":    response,
	})
}

// CancelReservation handles DELETE /api/v1/reservations/:id (Authenticated)
func (h *Handler) CancelReservation(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uuid.UUID)
	userRole, okRole := c.Get("user_role").(string)
	if !ok || !okRole {
		return c.JSON(http.StatusUnauthorized, map[string]any{
			"success": false,
			"message": "User not authenticated",
			"errors":  "Unauthorized",
		})
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "Invalid reservation ID format",
			"errors":  "ID must be a valid UUID",
		})
	}

	err = h.service.CancelReservation(userID, userRole, id)
	if err != nil {
		if errors.Is(err, ErrReservationNotFound) {
			return c.JSON(http.StatusNotFound, map[string]any{
				"success": false,
				"message": "Reservation not found",
				"errors":  err.Error(),
			})
		}
		if errors.Is(err, ErrForbiddenCancel) {
			return c.JSON(http.StatusForbidden, map[string]any{
				"success": false,
				"message": "Access forbidden: insufficient permissions",
				"errors":  err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to cancel reservation",
			"errors":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "Reservation cancelled successfully",
	})
}

// GetAllReservations handles GET /api/v1/reservations (Admin Only)
func (h *Handler) GetAllReservations(c *echo.Context) error {
	response, err := h.service.GetAllReservations()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to retrieve all reservations",
			"errors":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "Reservations retrieved successfully",
		"data":    response,
	})
}
