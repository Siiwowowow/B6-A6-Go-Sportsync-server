// internal/domin/reservation/service.go
package reservation

import (
	"errors"
	"gotickets/internal/domin/reservation/dto"

	"github.com/google/uuid"
)

var (
	ErrForbiddenCancel = errors.New("you are not authorized to cancel this reservation")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateReservation(userID uuid.UUID, req dto.CreateRequest) (*dto.Response, error) {
	res, err := s.repo.CreateWithLock(userID, req.ZoneID, req.LicensePlate)
	if err != nil {
		return nil, err
	}

	return &dto.Response{
		ID:           res.ID,
		UserID:       res.UserID,
		ZoneID:       res.ZoneID,
		LicensePlate: res.LicensePlate,
		Status:       res.Status,
		CreatedAt:    res.CreatedAt,
		UpdatedAt:    res.UpdatedAt,
	}, nil
}

func (s *Service) GetMyReservations(userID uuid.UUID) ([]*dto.Response, error) {
	reservations, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.Response, len(reservations))
	for i, res := range reservations {
		responses[i] = &dto.Response{
			ID:           res.ID,
			LicensePlate: res.LicensePlate,
			Status:       res.Status,
			Zone: &dto.ZoneInfo{
				ID:   res.Zone.ID,
				Name: res.Zone.Name,
				Type: res.Zone.Type,
			},
			CreatedAt: res.CreatedAt,
		}
	}

	return responses, nil
}

func (s *Service) CancelReservation(userID uuid.UUID, role string, id uuid.UUID) error {
	res, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// Drivers can only cancel their own reservations. Admins can cancel any.
	if role != "admin" && res.UserID != userID {
		return ErrForbiddenCancel
	}

	return s.repo.UpdateStatus(id, "cancelled")
}

func (s *Service) GetAllReservations() ([]*dto.Response, error) {
	reservations, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.Response, len(reservations))
	for i, res := range reservations {
		responses[i] = &dto.Response{
			ID:           res.ID,
			UserID:       res.UserID,
			ZoneID:       res.ZoneID,
			LicensePlate: res.LicensePlate,
			Status:       res.Status,
			CreatedAt:    res.CreatedAt,
			UpdatedAt:    res.UpdatedAt,
			Zone: &dto.ZoneInfo{
				ID:   res.Zone.ID,
				Name: res.Zone.Name,
				Type: res.Zone.Type,
			},
			User: &dto.UserInfo{
				ID:    res.User.ID,
				Name:  res.User.Name,
				Email: res.User.Email,
				Role:  res.User.Role,
			},
		}
	}

	return responses, nil
}
