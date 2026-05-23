package service

import (
	"context"

	"github.com/andrebarone77/cardiaflow-api/internal/domain"
	servicedto "github.com/andrebarone77/cardiaflow-api/internal/service/dto"
)

type HealthRecordRepository interface {
	Create(ctx context.Context, healthRecord *domain.HealthRecord) (string, error)
	GetByID(ctx context.Context, id string) (*domain.HealthRecord, error)
	ListByUserID(ctx context.Context, userId string) ([]*domain.HealthRecord, error)
	Update(ctx context.Context, id string, healthRecord *domain.HealthRecord) error
	Delete(ctx context.Context, id string) error
}

func NewHealthRecordService(repo HealthRecordRepository) *HealthRecordService {
	return &HealthRecordService{
		repo: repo,
	}
}

type HealthRecordService struct {
	repo HealthRecordRepository
}

func (s *HealthRecordService) Create(ctx context.Context, healthRecordInput servicedto.HealthRecordCreateInput) (string, error) {

	healthRecord := &domain.HealthRecord{
		UserID:             *healthRecordInput.UserID,
		HealthRecordTypeID: *healthRecordInput.HealthRecordTypeID,
		Value:              *healthRecordInput.Value,
		Notes:              &healthRecordInput.Notes,
		RecordedAt:         healthRecordInput.RecordedAt,
	}

	id, err := s.repo.Create(ctx, healthRecord)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *HealthRecordService) GetByID(ctx context.Context, id string) (*domain.HealthRecord, error) {

	healthRecord, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return healthRecord, nil

}

func (s *HealthRecordService) Update(ctx context.Context, id string, update servicedto.HealthRecordUpdateInput) error {

	healthRecord, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if update.Value != nil {
		healthRecord.Value = *update.Value
	}
	if update.Notes != nil {
		healthRecord.Notes = update.Notes
	}
	if update.RecordedAt != nil {
		healthRecord.RecordedAt = *update.RecordedAt
	}

	err = s.repo.Update(ctx, healthRecord.ID, healthRecord)

	return err
}

func (s *HealthRecordService) ListByUserID(ctx context.Context, userId string) ([]*domain.HealthRecord, error) {
	healthRecords, err := s.repo.ListByUserID(ctx, userId)

	if err != nil {
		return nil, err
	}

	return healthRecords, nil
}

func (s *HealthRecordService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
