package service

import (
	"context"
	"regexp"

	"github.com/andrebarone77/cardiaflow-api/internal/domain"
	servicedto "github.com/andrebarone77/cardiaflow-api/internal/service/dto"
)

type HealthRecordTypeRepository interface {
	Create(ctx context.Context, h *domain.HealthRecordType) (string, error)
	GetById(ctx context.Context, id string) (*domain.HealthRecordType, error)
	GetByCode(ctx context.Context, code string) (*domain.HealthRecordType, error)
	GetAll(ctx context.Context) ([]*domain.HealthRecordType, error)
	Update(ctx context.Context, id string, h *domain.HealthRecordType) error
	Delete(ctx context.Context, id string) error
	IsSystem(ctx context.Context, id string) (bool, error)
}

var codeRegex = regexp.MustCompile(`^[a-z][a-z0-9]*(?:_[a-z0-9]+)*$`)

func NewHealthRecordTypeService(repo HealthRecordTypeRepository) *HealthRecordTypeService {
	return &HealthRecordTypeService{
		repo: repo,
	}
}

type HealthRecordTypeService struct {
	repo HealthRecordTypeRepository
}

func (s *HealthRecordTypeService) Update(ctx context.Context, id string, update servicedto.HealthRecordTypeUpdateInput) (*domain.HealthRecordType, error) {

	isSystem, err := s.repo.IsSystem(ctx, id)
	if err != nil {
		return nil, err
	}
	if isSystem {
		return nil, domain.ErrHealthRecordTypeImmutable
	}

	healthTypeRecord, err := s.repo.GetById(ctx, id)

	if err != nil {
		return nil, err

	}

	if update.Code == nil && update.Name == nil && update.Unit == nil {
		return healthTypeRecord, domain.ErrNoInformation
	}

	if update.Code != nil {
		healthTypeRecord.Code = *update.Code
	}

	if update.Name != nil {
		healthTypeRecord.Name = *update.Name
	}

	if update.Unit != nil {
		healthTypeRecord.Unit = update.Unit
	}
	err = s.repo.Update(ctx, id, healthTypeRecord)

	return healthTypeRecord, nil
}

func (s *HealthRecordTypeService) Create(ctx context.Context, input servicedto.HealthRecordTypeInput) (string, error) {
	err := validateErrorCode(input.Code)

	if err != nil {
		return "", err
	}

	healthRecordType := &domain.HealthRecordType{
		Name: input.Name,
		Code: input.Code,
		Unit: input.Unit,
	}

	id, err := s.repo.Create(ctx, healthRecordType)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *HealthRecordTypeService) GetByID(ctx context.Context, id string) (*domain.HealthRecordType, error) {

	healthRecordType, err := s.repo.GetById(ctx, id)

	if err != nil {
		return nil, err
	}

	return healthRecordType, nil
}

func (s *HealthRecordTypeService) GetByCode(ctx context.Context, code string) (*domain.HealthRecordType, error) {

	healthRecordType, err := s.repo.GetByCode(ctx, code)

	if err != nil {
		return nil, err
	}

	return healthRecordType, nil
}

func (s *HealthRecordTypeService) GetAll(ctx context.Context) ([]*domain.HealthRecordType, error) {
	healthRecordTypes, err := s.repo.GetAll(ctx)

	if err != nil {
		return nil, err
	}

	return healthRecordTypes, nil

}

func (s *HealthRecordTypeService) Delete(ctx context.Context, id string) error {
	isSystem, err := s.repo.IsSystem(ctx, id)
	if err != nil {
		return err
	}
	if isSystem {
		return domain.ErrHealthRecordTypeImmutable
	}
	return s.repo.Delete(ctx, id)
}

func validateErrorCode(code string) error {
	if code == "" {
		return domain.ErrCodeRequired
	}

	if len(code) > 50 {
		return domain.ErrCodeTooLong
	}

	if !codeRegex.MatchString(code) {
		return domain.ErrCodeInvalid
	}

	return nil
}
