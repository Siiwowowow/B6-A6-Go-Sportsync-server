// internal/domin/zone/entity.go
package zone

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ParkingZone struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Name          string    `json:"name" gorm:"type:varchar(150);not null"`
	Type          string    `json:"type" gorm:"type:varchar(50);not null"` // general, ev_charging, covered
	TotalCapacity int       `json:"total_capacity" gorm:"not null"`
	PricePerHour  float64   `json:"price_per_hour" gorm:"type:decimal(10,2);not null"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (pz *ParkingZone) BeforeCreate(tx *gorm.DB) (err error) {
	pz.ID = uuid.New()
	return
}
