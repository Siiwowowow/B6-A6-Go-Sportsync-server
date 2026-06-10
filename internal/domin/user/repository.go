// internal/user/repository.go
package user

import (
	"errors"

	"gorm.io/gorm"
)

var ErrorAlreadyExists = errors.New("user already exists")

type Repository interface {
	CreateUser(user *User) error
	GetUserEmail(email string) (*User, error)
	GetUserByEmail(email string) (*User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}
func (r *repository) CreateUser(user *User) error {
	result := r.db.Create(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return ErrorAlreadyExists
		}
		return result.Error
	}
	return nil
}
func (r *repository) GetUserEmail(email string) (*User, error) {
	var user User
	result := r.db.Where(&User{Email: email}).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *repository) GetUserByEmail(email string) (*User, error) {
	return r.GetUserEmail(email)
}
