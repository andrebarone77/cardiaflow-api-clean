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

func (s *UserService) GetByEmail(ctx context.Context, email string, userID string, requesterRole domain.Role) (*domain.User, error) {

	user, err := s.repo.GetByEmail(ctx, email)

	if err != nil {
		return nil, err
	}

	if requesterRole == domain.RoleUser && user.ID != userID {
		return nil, domain.ErrForbidden
	}

	return user, nil

}

func (s *UserService) GetById(ctx context.Context, id string, userID string, requesterRole domain.Role) (*domain.User, error) {
	if requesterRole == domain.RoleUser && id != userID {
		return nil, domain.ErrForbidden
	}
	user, err := s.repo.GetById(ctx, id)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Delete(ctx context.Context, id string, userID string, requesterRole domain.Role) error {

	if requesterRole != domain.RoleAdmin && requesterRole != domain.RoleManager {

		return domain.ErrForbidden
	}

	if userID == id {
		return domain.ErrForbidden
	}

	return s.repo.Delete(ctx, id)

}

func (s *UserService) Update(ctx context.Context, id string, userId string, requesterRole domain.Role, req servicedto.UpdateUserInput) (*domain.User, error) {

	if requesterRole == domain.RoleUser && userId != id {
		return nil, domain.ErrForbidden
	}
	user, err := s.GetById(ctx, id, userId, requesterRole)

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
