package service

import (
	"context"
	"testing"
	"time"

	"github.com/andrebarone77/cardiaflow-api/internal/domain"
	servicedto "github.com/andrebarone77/cardiaflow-api/internal/service/dto"
)

type MockHealthRecordRepository struct {
	CreateFn       func(ctx context.Context, healthRecord *domain.HealthRecord) (string, error)
	GetByIDFn      func(ctx context.Context, id string) (*domain.HealthRecord, error)
	ListByUserIDFn func(ctx context.Context, userId string) ([]*domain.HealthRecord, error)
	UpdateFn       func(ctx context.Context, id string, healthRecord *domain.HealthRecord) error
	DeleteFn       func(ctx context.Context, id string) error
}

func (m *MockHealthRecordRepository) Create(ctx context.Context, healthRecord *domain.HealthRecord) (string, error) {
	return m.CreateFn(ctx, healthRecord)
}

func (m *MockHealthRecordRepository) GetByID(ctx context.Context, id string) (*domain.HealthRecord, error) {
	return m.GetByIDFn(ctx, id)
}

func (m *MockHealthRecordRepository) ListByUserID(ctx context.Context, userId string) ([]*domain.HealthRecord, error) {
	return m.ListByUserIDFn(ctx, userId)
}

func (m *MockHealthRecordRepository) Update(ctx context.Context, id string, healthRecord *domain.HealthRecord) error {
	return m.UpdateFn(ctx, id, healthRecord)
}

func (m *MockHealthRecordRepository) Delete(ctx context.Context, id string) error {
	return m.DeleteFn(ctx, id)
}

func TestHealthRecord_Create(t *testing.T) {
	noteOk := "TESTE OK"
	userIdOk := "22656871-b190-44b7-a313-b33d4d442d9b"
	healthRecordTypeOk := "94e9b6af-9d8d-43e5-8245-9c52cd4c435f"
	userIdInvalid := "03e8d5f0-be08-4ee5-ba0f-d1b3025f8fc4"
	valueOk := 3.5
	mockRepo := &MockHealthRecordRepository{
		CreateFn: func(ctx context.Context, h *domain.HealthRecord) (string, error) {
			if h.UserID == userIdInvalid {
				return "", domain.ErrUserNotFound
			}
			return "", nil
		},
	}

	service := NewHealthRecordService(mockRepo)

	tests := []struct {
		name         string
		healthRecord servicedto.HealthRecordCreateInput
		expect_error bool
	}{
		{
			name: "Test ok",
			healthRecord: servicedto.HealthRecordCreateInput{
				UserID:             &userIdOk,
				HealthRecordTypeID: &healthRecordTypeOk,
				Value:              &valueOk,
				Notes:              noteOk,
			},

			expect_error: false,
		},
		{
			name: "Misses UserID",
			healthRecord: servicedto.HealthRecordCreateInput{
				HealthRecordTypeID: &healthRecordTypeOk,
				Value:              &valueOk,
				Notes:              noteOk,
			},
			expect_error: true,
		},
		{
			name: "Misses Note",
			healthRecord: servicedto.HealthRecordCreateInput{
				UserID:             &userIdOk,
				HealthRecordTypeID: &healthRecordTypeOk,
				Value:              &valueOk,
			},
			expect_error: true,
		}, {
			name: "Test - User Not Found",
			healthRecord: servicedto.HealthRecordCreateInput{
				UserID:             &userIdInvalid,
				HealthRecordTypeID: &healthRecordTypeOk,
				Value:              &valueOk,
				Notes:              noteOk,
			},

			expect_error: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.Create(context.Background(), tt.healthRecord)
			if tt.expect_error && err == nil {
				t.Errorf("Expected error")
			}
			if !tt.expect_error && err != nil {
				t.Errorf("Error not expected")
			}
		})
	}
}

func TestHealthRecord_GetById(t *testing.T) {
	idOk := "3d74927e-0604-4f64-823c-e66a0c84da60"
	idNotFound := "d8b5f97a-a06b-4d1f-8953-ae478e4ee835"
	mockRepo := &MockHealthRecordRepository{
		GetByIDFn: func(ctx context.Context, id string) (*domain.HealthRecord, error) {
			if id == idNotFound {
				return nil, domain.ErrCodeInvalid
			}
			return nil, nil
		},
	}

	service := NewHealthRecordService(mockRepo)

	tests := []struct {
		name         string
		id           string
		expect_error bool
	}{
		{
			name:         "Test Ok",
			id:           idOk,
			expect_error: false,
		}, {
			name:         "ID not found",
			id:           idNotFound,
			expect_error: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.GetByID(context.Background(), tt.id)
			if tt.expect_error && err == nil {
				t.Errorf("Expected error")
			}
			if !tt.expect_error && err != nil {
				t.Errorf("Error not expected")
			}
		})
	}

}

func TestHealthRecord_Update(t *testing.T) {
	idOk := "3d74927e-0604-4f64-823c-e66a0c84da60"
	idNotFound := "d8b5f97a-a06b-4d1f-8953-ae478e4ee835"
	userIdOk := "573f14b1-e1a4-4b59-a5cb-e85f308be79b"
	valueOk := float64(0.0)
	notesOK := "Notes OK"
	healthRecordTypId := "helath_record_typ_id"
	recordedAt := time.Now()
	mockRepo := &MockHealthRecordRepository{
		GetByIDFn: func(ctx context.Context, id string) (*domain.HealthRecord, error) {
			if id == idNotFound {
				return nil, domain.ErrCodeInvalid
			}
			return &domain.HealthRecord{
				ID:                 idOk,
				UserID:             userIdOk,
				HealthRecordTypeID: healthRecordTypId,
				Value:              valueOk,
				RecordedAt:         time.Now(),
				Notes:              &notesOK,
				CreatedAt:          time.Now().Add(1 * time.Minute),
				UpdatedAt:          time.Now().Add(3 * time.Minute),
			}, nil
		},
		UpdateFn: func(ctx context.Context, id string, healthRecord *domain.HealthRecord) error {
			return nil
		},
	}

	service := NewHealthRecordService(mockRepo)

	tests := []struct {
		name         string
		id           string
		update_input servicedto.HealthRecordUpdateInput
		expect_error bool
	}{
		{
			name: "Test OK",
			id:   idOk,
			update_input: servicedto.HealthRecordUpdateInput{
				Value:      &valueOk,
				Notes:      &notesOK,
				RecordedAt: &recordedAt,
			},
			expect_error: false,
		},
		{
			name: "Test - ID not found",
			id:   idNotFound,
			update_input: servicedto.HealthRecordUpdateInput{
				Value:      &valueOk,
				Notes:      &notesOK,
				RecordedAt: &recordedAt,
			},
			expect_error: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.Update(context.Background(), tt.id, tt.update_input)
			if tt.expect_error && err == nil {
				t.Errorf("Expected error")
			}
			if !tt.expect_error && err != nil {
				t.Errorf("Error not expected")
			}
		})
	}
}

func TestHealthRecord_ListByUserID(t *testing.T) {

	userIdOk := "573f14b1-e1a4-4b59-a5cb-e85f308be79b"
	userNotFound := "1774ac5a-58e0-4b26-abe8-f1bd0e5fde0b"
	mockRepo := &MockHealthRecordRepository{
		ListByUserIDFn: func(ctx context.Context, userId string) ([]*domain.HealthRecord, error) {
			if userId == userNotFound {
				return nil, domain.ErrUserNotFound
			}
			return nil, nil
		},
	}

	service := NewHealthRecordService(mockRepo)

	tests := []struct {
		name         string
		userId       string
		expect_error bool
	}{
		{
			name:         "Test OK",
			userId:       userIdOk,
			expect_error: false,
		},
		{
			name:         "User not found",
			userId:       userNotFound,
			expect_error: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.ListByUserID(context.Background(), tt.userId)

			if tt.expect_error && err == nil {
				t.Errorf("Error expected")
			}

			if !tt.expect_error && err != nil {
				t.Errorf("Unexpcted error")
			}
		})
	}
}

func TestHealthRecord_Delete(t *testing.T) {
	idOk := "b9f56f1b-2119-40d5-b442-068c6f75c55e"
	idNotFound := "9cd74b61-5b94-4d04-a748-4c731ef8657a"
	mockRepo := &MockHealthRecordRepository{
		DeleteFn: func(ctx context.Context, id string) error {
			if id == idNotFound {
				return domain.ErrHealthRecordNotFound
			}
			return nil
		},
	}

	service := NewHealthRecordService(mockRepo)

	tests := []struct {
		name         string
		id           string
		expect_error bool
	}{
		{
			name:         "Test OK",
			id:           idOk,
			expect_error: false,
		},
		{
			name:         "Test OK",
			id:           idNotFound,
			expect_error: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.Delete(context.Background(), tt.id)

			if tt.expect_error && err == nil {
				t.Errorf("Error expected")
			}

			if !tt.expect_error && err != nil {
				t.Errorf("Unexpcted error")
			}
		})
	}

}
