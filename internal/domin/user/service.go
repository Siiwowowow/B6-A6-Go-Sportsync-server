// internal/user/service.go
package user

import (
	"fmt"
	"gotickets/internal/domin/user/dto"
	"gotickets/internal/auth"
)

var ErrInvalidCredentials = fmt.Errorf("invalid email or password")

type Service struct {
	repo Repository
	jwtService auth.JWTService
}

func NewService(repo Repository, jwtService auth.JWTService) *Service {
	return &Service{
		repo: repo,
		jwtService: jwtService,
	}
}
func (s *Service) CreateUser(req dto.CreateRequest) (*dto.Response, error) {
	user := &User{
		Name:  req.Name,
		Email: req.Email,
	}
	err := user.hashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	//generate JWT token
	
	err = s.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}
	response := &dto.Response{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
	}
	return response, nil

}
func (s *Service) LoginUser(req dto.LoginRequest) (*dto.Response, error) {
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
	token, err := s.jwtService.GenerateToken(user.ID, user.Name, user.Email)
	if err != nil {
		return nil,fmt.Errorf("failed to generate token: %w", err)
	}
	response := &dto.Response{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Token:     token,
		CreatedAt: user.CreatedAt.String(),
	}
	return response, nil

}
