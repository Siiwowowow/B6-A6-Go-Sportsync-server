// internal/domin/zone/handler.go
package zone

import (
	"errors"
	"gotickets/internal/domin/zone/dto"
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

// CreateZone handles POST /api/v1/zones (Admin Only)
func (h *Handler) CreateZone(c *echo.Context) error {
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

	response, err := h.service.CreateZone(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to create parking zone",
			"errors":  err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]any{
		"success": true,
		"message": "Parking zone created successfully",
		"data":    response,
	})
}

// GetAllZones handles GET /api/v1/zones (Public)
func (h *Handler) GetAllZones(c *echo.Context) error {
	response, err := h.service.GetAllZones()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to retrieve parking zones",
			"errors":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "Parking zones retrieved successfully",
		"data":    response,
	})
}

// GetZoneByID handles GET /api/v1/zones/:id (Public)
func (h *Handler) GetZoneByID(c *echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"success": false,
			"message": "Invalid zone ID format",
			"errors":  "ID must be a valid UUID",
		})
	}

	response, err := h.service.GetZoneByID(id)
	if err != nil {
		if errors.Is(err, ErrZoneNotFound) {
			return c.JSON(http.StatusNotFound, map[string]any{
				"success": false,
				"message": "Parking zone not found",
				"errors":  err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Failed to retrieve parking zone",
			"errors":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"message": "Parking zone retrieved successfully",
		"data":    response,
	})
}
