// internal/domin/zone/repository.go
package zone

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrZoneNotFound = errors.New("parking zone not found")
)

type Repository interface {
	Create(zone *ParkingZone) error
	FindByID(id uuid.UUID) (*ParkingZone, error)
	FindAll() ([]*ParkingZone, error)
	GetActiveReservationsCount(zoneID uuid.UUID) (int, error)
	GetActiveReservationsCounts() (map[uuid.UUID]int, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(zone *ParkingZone) error {
	return r.db.Create(zone).Error
}

func (r *repository) FindByID(id uuid.UUID) (*ParkingZone, error) {
	var zone ParkingZone
	err := r.db.First(&zone, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrZoneNotFound
		}
		return nil, err
	}
	return &zone, nil
}

func (r *repository) FindAll() ([]*ParkingZone, error) {
	var zones []*ParkingZone
	err := r.db.Find(&zones).Error
	if err != nil {
		return nil, err
	}
	return zones, nil
}

func (r *repository) GetActiveReservationsCount(zoneID uuid.UUID) (int, error) {
	var count int64
	err := r.db.Table("reservations").
		Where("zone_id = ? AND status = ?", zoneID, "active").
		Count(&count).Error
	return int(count), err
}

func (r *repository) GetActiveReservationsCounts() (map[uuid.UUID]int, error) {
	type Result struct {
		ZoneID uuid.UUID `gorm:"column:zone_id"`
		Count  int       `gorm:"column:count"`
	}
	var results []Result
	err := r.db.Table("reservations").
		Select("zone_id, count(*) as count").
		Where("status = ?", "active").
		Group("zone_id").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	counts := make(map[uuid.UUID]int)
	for _, res := range results {
		counts[res.ZoneID] = res.Count
	}
	return counts, nil
}
