// internal/domin/reservation/dto/response.go
package dto

import (
	"time"

	"github.com/google/uuid"
)

type ZoneInfo struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Type string    `json:"type"`
}

type UserInfo struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Role  string    `json:"role"`
}

type Response struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id,omitempty"`
	User         *UserInfo `json:"user,omitempty"`
	ZoneID       uuid.UUID `json:"zone_id"`
	Zone         *ZoneInfo `json:"zone,omitempty"`
	LicensePlate string    `json:"license_plate"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}
