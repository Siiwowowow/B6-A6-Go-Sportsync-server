// internal/user/service.go
package user

import (
	"fmt"
	"gotickets/internal/auth"
	"gotickets/internal/domin/user/dto"
)

var ErrInvalidCredentials = fmt.Errorf("invalid email or password")

type Service struct {
	repo       Repository
	jwtService auth.JWTService
}

func NewService(repo Repository, jwtService auth.JWTService) *Service {
	return &Service{
		repo:       repo,
		jwtService: jwtService,
	}
}

func (s *Service) CreateUser(req dto.CreateRequest) (*dto.UserResponse, error) {
	role := req.Role
	if role == "" {
		role = "driver"
	}

	user := &User{
		Name:  req.Name,
		Email: req.Email,
		Role:  role,
	}
	err := user.hashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	
	err = s.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	response := &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	return response, nil
}

func (s *Service) LoginUser(req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrInvalidCredentials
	}
	err = user.checkPassword(req.Password)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	token, err := s.jwtService.GenerateToken(user.ID, user.Name, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	response := &dto.LoginResponse{
		Token: token,
		User: dto.UserMini{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	}
	return response, nil
}
