package service

import (
	"context"
	"log"
	"strings"

	"github.com/andrebarone77/cardiaflow-api/internal/domain"
	servicedto "github.com/andrebarone77/cardiaflow-api/internal/service/dto"
	"github.com/andrebarone77/cardiaflow-api/pkg/utils"
)

type UserRepository interface {
	Save(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetById(ctx context.Context, id string) (*domain.User, error)
	Delete(ctx context.Context, id string) error
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, req servicedto.CreateUserInput) (*domain.User, error) {

	hash, err := utils.HashPassword(req.Password)

	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Name:         req.Name,
		Email:        strings.ToLower(strings.TrimSpace(req.Email)),
		PasswordHash: hash,
	}

	if err = s.repo.Save(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {

	user, err := s.repo.GetByEmail(ctx, email)

	if err != nil {
		return nil, err
	}

	return user, nil

}

func (s *UserService) GetById(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.repo.GetById(ctx, id)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)

}

func (s *UserService) Update(ctx context.Context, id string, req servicedto.UpdateUserInput) (*domain.User, error) {

	user, err := s.GetById(ctx, id)

	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		user.Name = *req.Name
	}

	if req.Email != nil {
		user.Email = strings.ToLower(strings.TrimSpace(*req.Email))
	}

	if req.Password != nil {
		hash, err := utils.HashPassword(*req.Password)

		if err != nil {
			return nil, err
		}
		user.PasswordHash = hash
	}

	log.Printf("user.password_hash[%s]", user.PasswordHash)

	err = s.repo.Save(ctx, user)

	if err != nil {
		return nil, err
	}

	return user, nil
}
