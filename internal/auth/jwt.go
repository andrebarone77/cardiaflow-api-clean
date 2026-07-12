package auth

import (
	"errors"
	"time"

	"github.com/andrebarone77/cardiaflow-api/configs"
	"github.com/andrebarone77/cardiaflow-api/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Role domain.Role `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(userID string, role domain.Role) (string, error) {
	cfg := configs.Load()

	now := time.Now()

	duration, err := time.ParseDuration(cfg.JWTExpiresIn)
	if err != nil {
		return "", err
	}

	claims := Claims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(cfg.JWTSecret))
}

func ValidateToken(tokenString string) (*Claims, error) {
	cfg := configs.Load()

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, errors.New("unexpected signing method")
			}

			return []byte(cfg.JWTSecret), nil
		})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func (c *Claims) UserID() string {
	return c.Subject
}
