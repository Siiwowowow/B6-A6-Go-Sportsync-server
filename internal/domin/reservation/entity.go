// internal/domin/reservation/entity.go
package reservation

import (
	"gotickets/internal/domin/user"
	"gotickets/internal/domin/zone"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Reservation struct {
	ID           uuid.UUID        `json:"id" gorm:"type:uuid;primaryKey"`
	UserID       uuid.UUID        `json:"user_id" gorm:"type:uuid;not null"`
	User         user.User        `json:"user" gorm:"foreignKey:UserID"`
	ZoneID       uuid.UUID        `json:"zone_id" gorm:"type:uuid;not null"`
	Zone         zone.ParkingZone `json:"zone" gorm:"foreignKey:ZoneID"`
	LicensePlate string           `json:"license_plate" gorm:"type:varchar(15);not null"`
	Status       string           `json:"status" gorm:"type:varchar(20);default:'active';not null"` // active, completed, cancelled
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

func (r *Reservation) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.New()
	return
}
