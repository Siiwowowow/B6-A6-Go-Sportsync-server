// internal/auth/jwt.go
package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	jwtSecretKey         = "your_secret_key"
	defaultTokenDuration = 24 * time.Hour
)

type JWTClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Name   string    `json:"name"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

type JWTService interface {
	GenerateToken(userId uuid.UUID, name string, email string, role string) (string, error)
	ValidateToken(tokenStr string) (*JWTClaims, error)
}

type jwtService struct {
	secretKey string
	duration  time.Duration
}

func NewJWTService(secretKey string, tokenDuration time.Duration) JWTService {
	if secretKey == "" {
		secretKey = jwtSecretKey
	}
	if tokenDuration == 0 {
		tokenDuration = defaultTokenDuration
	}
	return &jwtService{
		secretKey: secretKey,
		duration:  tokenDuration,
	}
}

func (js *jwtService) GenerateToken(userId uuid.UUID, name string, email string, role string) (string, error) {
	claims := &JWTClaims{
		UserID: userId,
		Name:   name,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(js.duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "gotickets",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(js.secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (js *jwtService) ValidateToken(tokenStr string) (*JWTClaims, error) {
	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method: " + token.Header["alg"].(string))
		}
		return []byte(js.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	tokenClaims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	return tokenClaims, nil
}
