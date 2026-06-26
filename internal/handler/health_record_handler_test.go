package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/andrebarone77/cardiaflow-api/internal/domain"
	servicedto "github.com/andrebarone77/cardiaflow-api/internal/service/dto"
	"github.com/gin-gonic/gin"
)

type MockHealthRecordService struct {
	CreateFn       func(ctx context.Context, healthRecordInput servicedto.HealthRecordCreateInput) (string, error)
	GetByIDFn      func(ctx context.Context, id string) (*domain.HealthRecord, error)
	UpdateFn       func(ctx context.Context, id string, update_input servicedto.HealthRecordUpdateInput) error
	ListByUserIDFn func(ctx context.Context, userId string) ([]*domain.HealthRecord, error)
	DeleteFn       func(ctx context.Context, id string) error
}

func (m *MockHealthRecordService) Create(ctx context.Context, healthRecordInput servicedto.HealthRecordCreateInput) (string, error) {
	return m.CreateFn(ctx, healthRecordInput)
}

func (m *MockHealthRecordService) GetByID(ctx context.Context, id string) (*domain.HealthRecord, error) {
	return m.GetByIDFn(ctx, id)
}

func (m *MockHealthRecordService) Update(ctx context.Context, id string, update_input servicedto.HealthRecordUpdateInput) error {
	return m.UpdateFn(ctx, id, update_input)
}

func (m *MockHealthRecordService) ListByUserID(ctx context.Context, userId string) ([]*domain.HealthRecord, error) {
	return m.ListByUserIDFn(ctx, userId)
}

func (m *MockHealthRecordService) Delete(ctx context.Context, id string) error {
	return m.DeleteFn(ctx, id)
}

func TestHealthRecord_Create(t *testing.T) {
	errNum := 0
	mockHealthRecordService := &MockHealthRecordService{
		CreateFn: func(ctx context.Context, healthRecordInput servicedto.HealthRecordCreateInput) (string, error) {
			switch errNum {
			case 1:
				errNum = 0
				return "", domain.ErrInvalidUserOrHealthRecordType
			case 2:
				errNum = 0
				return "", errors.New("Generic Errro")
			}
			return "", nil
		},
	}
	handler := NewHealthRecordHandler(mockHealthRecordService)

	tests := []struct {
		test_name       string
		body            string
		expected_status int
		err_num         int
	}{
		{
			test_name: "Test OK",
			body: `
			{
				"user_id": "userID",
				"health_record_type_id": "type",
				"value": 3.0,
				"notes": "Notes Ok",
				"recorded_at": "2026-06-25T19:30:00Z"
			}
			`,
			expected_status: http.StatusOK,
			err_num:         0,
		},
		{
			test_name: "Test Invalid ID",
			body: `
			{
				"user_id": "userID",
				"health_record_type_id": "invalid_type",
				"value": 3.0,
				"notes": "Notes Ok",
				"recorded_at": "2026-06-25T19:30:00Z"
			}
			`,
			expected_status: http.StatusBadRequest,
			err_num:         1,
		},
		{
			test_name: "Test Generic Error",
			body: `
			{
				"user_id": "userID",
				"health_record_type_id": "invalid_type",
				"value": 3.0,
				"notes": "Notes Ok",
				"recorded_at": "2026-06-25T19:30:00Z"
			}
			`,
			expected_status: http.StatusInternalServerError,
			err_num:         2,
		},
		{
			test_name: "Test Invalid Body",
			body: `
			{
				"user_id": "userID",
				"health_record_type_id": "invalid_type",
				"value": 3.0,
				"notes": "Notes Ok,
				"recorded_at": "2026-06-25T19:30:00Z"
			}
			`,
			expected_status: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.test_name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(
				http.MethodPost,
				"/api/healthrecord",
				strings.NewReader(tt.body),
			)
			errNum = tt.err_num
			handler.Create(c)
			status := c.Writer.Status()

			if status != tt.expected_status {
				t.Errorf("Unexpected status %v", status)
			}

		})
	}

}

func TestHealthRecord_Update(t *testing.T) {
	mockHealthRecordService := &MockHealthRecordService{
		UpdateFn: func(ctx context.Context, id string, update_input servicedto.HealthRecordUpdateInput) error {
			if id == ID_NOT_FOUND {
				return domain.ErrHealthRecordNotFound
			}

			if id == ID_GENERIC_ERROR {
				return errors.New("Generic Error")
			}
			return nil
		},
	}

	handler := NewHealthRecordHandler(mockHealthRecordService)

	tests := []struct {
		test_name       string
		id              string
		body            string
		expected_status int
	}{
		{
			test_name: "Test Ok",
			id:        ID_OK,
			body: `
			{
				"value": 3.0,
				"notes": "Notes Ok",
				"recorded_at": "2026-06-25T19:30:00Z"
			}
			`,
			expected_status: http.StatusOK,
		},
		{
			test_name: "Test Missing ID",
			body: `
			{
				"value": 3.0,
				"notes": "Notes Ok",
				"recorded_at": "2026-06-25T19:30:00Z"
			}
			`,
			expected_status: http.StatusBadRequest,
		},
		{
			test_name: "Test Malformed Body",
			id:        ID_OK,
			body: `
			{
				"value": 3.0,
				"notes": "Notes Ok,
				"recorded_at": "2026-06-25T19:30:00Z"
			}
			`,
			expected_status: http.StatusBadRequest,
		},
		{
			test_name: "Test Not Found",
			id:        ID_NOT_FOUND,
			body: `
			{
				"value": 3.0,
				"notes": "Notes Ok",
				"recorded_at": "2026-06-25T19:30:00Z"
			}
			`,
			expected_status: http.StatusNotFound,
		},
		{
			test_name: "Test Generic Error",
			id:        ID_GENERIC_ERROR,
			body: `
			{
				"value": 3.0,
				"notes": "Notes Ok",
				"recorded_at": "2026-06-25T19:30:00Z"
			}
			`,
			expected_status: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.test_name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(
				http.MethodPatch,
				"/api/healthrecord",
				strings.NewReader(tt.body),
			)
			if tt.id != "" {
				c.Params = gin.Params{
					{
						Key:   "id",
						Value: tt.id,
					},
				}
			}
			handler.Update(c)
			status := c.Writer.Status()
			if status != tt.expected_status {
				t.Errorf("Unexpected Status %v", status)
			}
		})
	}

}
func TestHealhRecordType_GetById(t *testing.T) {
	mockHealthRecordService := &MockHealthRecordService{
		GetByIDFn: func(ctx context.Context, id string) (*domain.HealthRecord, error) {
			if id == ID_NOT_FOUND {
				return nil, domain.ErrHealthRecordNotFound
			}
			if id == ID_GENERIC_ERROR {
				return nil, errors.New("Generic Error")
			}
			return nil, nil
		},
	}

	handler := NewHealthRecordHandler(mockHealthRecordService)

	tests := []struct {
		test_name       string
		id              string
		expected_status int
	}{
		{
			test_name:       "Test OK",
			id:              ID_OK,
			expected_status: http.StatusOK,
		},
		{
			test_name:       "Test Not Found",
			id:              ID_NOT_FOUND,
			expected_status: http.StatusNotFound,
		},
		{
			test_name:       "Test Generic Error",
			id:              ID_GENERIC_ERROR,
			expected_status: http.StatusInternalServerError,
		},
		{
			test_name:       "Test missing Id",
			expected_status: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.test_name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(
				http.MethodGet,
				"/api/healthrecord",
				nil,
			)

			if tt.id != "" {
				c.Params = gin.Params{
					{
						Key:   "id",
						Value: tt.id,
					},
				}
			}
			handler.GetByID(c)
			status := c.Writer.Status()

			if status != tt.expected_status {
				t.Errorf("Unexpected Stauts %v", status)
			}

		})
	}
}

func TestHealthRecord_ListByUserID(t *testing.T) {

	mockHealthRecord := &MockHealthRecordService{
		ListByUserIDFn: func(ctx context.Context, userId string) ([]*domain.HealthRecord, error) {
			notes := "Notes Ok"

			if userId == ID_NOT_FOUND {
				return nil, domain.ErrHealthRecordNotFound
			}

			if userId == ID_EMPTY_RETURN {
				return []*domain.HealthRecord{}, nil
			}

			if userId == ID_GENERIC_ERROR {
				return nil, errors.New("Generic Error")
			}

			return []*domain.HealthRecord{
				{
					ID:                 ID_OK,
					UserID:             ID_OK,
					HealthRecordTypeID: ID_OK,
					Value:              3.0,
					Notes:              &notes,
					RecordedAt:         time.Now().Add(-3 * time.Minute),
					CreatedAt:          time.Now().Add(-2 * time.Minute),
					UpdatedAt:          time.Now().Add(time.Minute),
				},
				{
					ID:                 ID_OK,
					UserID:             ID_OK,
					HealthRecordTypeID: ID_OK,
					Value:              3.0,
					Notes:              &notes,
					RecordedAt:         time.Now().Add(-3 * time.Minute),
					CreatedAt:          time.Now().Add(-2 * time.Minute),
					UpdatedAt:          time.Now().Add(time.Minute),
				},
			}, nil
		},
	}

	handler := NewHealthRecordHandler(mockHealthRecord)

	tests := []struct {
		test_name       string
		user_id         string
		expected_status int
	}{
		{
			test_name:       "Test Ok",
			user_id:         ID_OK,
			expected_status: http.StatusOK,
		},
		{
			test_name:       "Test Not Found",
			user_id:         ID_NOT_FOUND,
			expected_status: http.StatusNotFound,
		},
		{
			test_name:       "Test Generic Error",
			user_id:         ID_GENERIC_ERROR,
			expected_status: http.StatusInternalServerError,
		},
		{
			test_name:       "Test Generic Error",
			user_id:         ID_EMPTY_RETURN,
			expected_status: http.StatusNotFound,
		},

		{
			test_name:       "Test Empty UserID",
			user_id:         "",
			expected_status: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.test_name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(
				http.MethodGet,
				fmt.Sprintf("/api/healthrecord?user_id=%v", tt.user_id),
				nil,
			)

			handler.ListByUserID(c)
			status := c.Writer.Status()

			if status != tt.expected_status {
				t.Errorf("Unexpected Status %v", status)
			}

		})
	}
}

func TestHandlerRecord_Delete(t *testing.T) {
	mockHealthRecordService := &MockHealthRecordService{
		DeleteFn: func(ctx context.Context, id string) error {
			if id == ID_NOT_FOUND {
				return domain.ErrHealthRecordNotFound
			}

			if id == ID_GENERIC_ERROR {
				return errors.New("Generic Error")
			}
			return nil
		},
	}

	handler := NewHealthRecordHandler(mockHealthRecordService)

	tests := []struct {
		test_name       string
		id              string
		expected_status int
	}{
		{
			test_name:       "Test Ok",
			id:              ID_OK,
			expected_status: http.StatusNoContent,
		},
		{
			test_name:       "Test Not Found",
			id:              ID_NOT_FOUND,
			expected_status: http.StatusNotFound,
		},
		{
			test_name:       "Test Generic Error",
			id:              ID_GENERIC_ERROR,
			expected_status: http.StatusInternalServerError,
		},

		{
			test_name:       "Test No ID",
			id:              "",
			expected_status: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.test_name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(
				http.MethodDelete,
				fmt.Sprintf("/api/healthrecord?id=%v", tt.id),
				nil,
			)

			handler.Delete(c)
			status := c.Writer.Status()

			if status != tt.expected_status {
				t.Errorf("Unexpected Status %v", status)
			}

		})
	}

}
