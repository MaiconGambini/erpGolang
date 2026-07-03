package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID   string `json:"user_id"`
	TenantID string `json:"tenant_id"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type Service struct {
	secret []byte
	ttl    time.Duration
}

func NewService(secret string, ttl time.Duration) Service {
	return Service{secret: []byte(secret), ttl: ttl}
}

func (s Service) GenerateAccessToken(userID, tenantID, role string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		TenantID: tenantID,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.ttl)),
			ID:        uuid.NewString(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s Service) Validate(tokenString string) (Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return s.secret, nil
	})
	if err != nil {
		return Claims{}, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return Claims{}, fmt.Errorf("invalid token")
	}
	return *claims, nil
}
