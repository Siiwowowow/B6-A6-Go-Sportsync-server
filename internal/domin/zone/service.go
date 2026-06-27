// internal/domin/zone/service.go
package zone

import (
	"gotickets/internal/domin/zone/dto"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateZone(req dto.CreateRequest) (*dto.Response, error) {
	zone := &ParkingZone{
		Name:          req.Name,
		Type:          req.Type,
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}

	if err := s.repo.Create(zone); err != nil {
		return nil, err
	}

	return &dto.Response{
		ID:             zone.ID,
		Name:           zone.Name,
		Type:           zone.Type,
		TotalCapacity:  zone.TotalCapacity,
		AvailableSpots: zone.TotalCapacity,
		PricePerHour:   zone.PricePerHour,
		CreatedAt:      zone.CreatedAt,
		UpdatedAt:      zone.UpdatedAt,
	}, nil
}

func (s *Service) GetZoneByID(id uuid.UUID) (*dto.Response, error) {
	zone, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	activeCount, err := s.repo.GetActiveReservationsCount(zone.ID)
	if err != nil {
		return nil, err
	}

	availableSpots := zone.TotalCapacity - activeCount
	if availableSpots < 0 {
		availableSpots = 0
	}

	return &dto.Response{
		ID:             zone.ID,
		Name:           zone.Name,
		Type:           zone.Type,
		TotalCapacity:  zone.TotalCapacity,
		AvailableSpots: availableSpots,
		PricePerHour:   zone.PricePerHour,
		CreatedAt:      zone.CreatedAt,
		UpdatedAt:      zone.UpdatedAt,
	}, nil
}

func (s *Service) GetAllZones() ([]*dto.Response, error) {
	zones, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	activeCounts, err := s.repo.GetActiveReservationsCounts()
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.Response, len(zones))
	for i, zone := range zones {
		activeCount := activeCounts[zone.ID]
		availableSpots := zone.TotalCapacity - activeCount
		if availableSpots < 0 {
			availableSpots = 0
		}

		responses[i] = &dto.Response{
			ID:             zone.ID,
			Name:           zone.Name,
			Type:           zone.Type,
			TotalCapacity:  zone.TotalCapacity,
			AvailableSpots: availableSpots,
			PricePerHour:   zone.PricePerHour,
			CreatedAt:      zone.CreatedAt,
		}
	}

	return responses, nil
}
