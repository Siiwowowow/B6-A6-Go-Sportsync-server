// internal/domin/reservation/repository.go
package reservation

import (
	"errors"
	"gotickets/internal/domin/zone"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrZoneFull            = errors.New("parking zone is full")
	ErrReservationNotFound = errors.New("reservation not found")
)

type Repository interface {
	CreateWithLock(userID uuid.UUID, zoneID uuid.UUID, licensePlate string) (*Reservation, error)
	GetByID(id uuid.UUID) (*Reservation, error)
	GetByUserID(userID uuid.UUID) ([]*Reservation, error)
	GetAll() ([]*Reservation, error)
	UpdateStatus(id uuid.UUID, status string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateWithLock(userID uuid.UUID, zoneID uuid.UUID, licensePlate string) (*Reservation, error) {
	var reservation Reservation

	err := r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Lock the parking zone row to prevent race conditions on capacity check
		var pZone zone.ParkingZone
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&pZone, zoneID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return zone.ErrZoneNotFound
			}
			return err
		}

		// 2. Count active reservations in this zone
		var activeCount int64
		if err := tx.Model(&Reservation{}).Where("zone_id = ? AND status = ?", zoneID, "active").Count(&activeCount).Error; err != nil {
			return err
		}

		// 3. Verify capacity
		if int(activeCount) >= pZone.TotalCapacity {
			return ErrZoneFull
		}

		// 4. Create reservation
		reservation = Reservation{
			UserID:       userID,
			ZoneID:       zoneID,
			LicensePlate: licensePlate,
			Status:       "active",
		}
		if err := tx.Create(&reservation).Error; err != nil {
			return err
		}

		// 5. Preload the relations so the service gets complete data
		if err := tx.Preload("Zone").Preload("User").First(&reservation, reservation.ID).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return &reservation, nil
}

func (r *repository) GetByID(id uuid.UUID) (*Reservation, error) {
	var res Reservation
	err := r.db.Preload("Zone").Preload("User").First(&res, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReservationNotFound
		}
		return nil, err
	}
	return &res, nil
}

func (r *repository) GetByUserID(userID uuid.UUID) ([]*Reservation, error) {
	var reservations []*Reservation
	err := r.db.Preload("Zone").Where("user_id = ?", userID).Find(&reservations).Error
	if err != nil {
		return nil, err
	}
	return reservations, nil
}

func (r *repository) GetAll() ([]*Reservation, error) {
	var reservations []*Reservation
	err := r.db.Preload("Zone").Preload("User").Find(&reservations).Error
	if err != nil {
		return nil, err
	}
	return reservations, nil
}

func (r *repository) UpdateStatus(id uuid.UUID, status string) error {
	result := r.db.Model(&Reservation{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrReservationNotFound
	}
	return nil
}
