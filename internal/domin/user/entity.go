// internal/user/entity.go
package user

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User represents the GORM schema and properties of a user in SpotSync.
type User struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Name      string         `json:"name" validate:"required" gorm:"type:varchar(100);not null"`
	Email     string         `json:"email" validate:"required,email" gorm:"type:varchar(100);uniqueIndex;not null"`
	Password  string         `json:"password" validate:"required,min=6" gorm:"type:varchar(255);not null"`
	Role      string         `json:"role" gorm:"type:varchar(20);default:'driver';not null"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}

func (u *User) hashPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

func (u *User) checkPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}