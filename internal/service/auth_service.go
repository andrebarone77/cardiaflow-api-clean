package service

import (
	"context"

	"github.com/andrebarone77/cardiaflow-api/internal/auth"
	"github.com/andrebarone77/cardiaflow-api/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepository interface {
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

type AuthService struct {
	repo AuthRepository
}

func NewAuthService(repo AuthRepository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (string, error) {

	user, err := s.repo.GetByEmail(ctx, email)

	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", domain.ErrNotAuthorized
	}

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
