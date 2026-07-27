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
	GetByIDFn      func(ctx context.Context, id string, userID string) (*domain.HealthRecord, error)
	UpdateFn       func(ctx context.Context, id string, update_input servicedto.HealthRecordUpdateInput, userID string) error
	ListByUserIDFn func(ctx context.Context, userId string, requesterID string) ([]*domain.HealthRecord, error)
	DeleteFn       func(ctx context.Context, id string, requesterID string) error
}

func (m *MockHealthRecordService) Create(ctx context.Context, healthRecordInput servicedto.HealthRecordCreateInput) (string, error) {
	return m.CreateFn(ctx, healthRecordInput)
}

func (m *MockHealthRecordService) GetByID(ctx context.Context, id string, userID string) (*domain.HealthRecord, error) {
	return m.GetByIDFn(ctx, id, userID)
}

func (m *MockHealthRecordService) Update(ctx context.Context, id string, update_input servicedto.HealthRecordUpdateInput, userID string) error {
	return m.UpdateFn(ctx, id, update_input, userID)
}

func (m *MockHealthRecordService) ListByUserID(ctx context.Context, userId string, requesterID string) ([]*domain.HealthRecord, error) {
	return m.ListByUserIDFn(ctx, userId, requesterID)
}

func (m *MockHealthRecordService) Delete(ctx context.Context, id string, requesterID string) error {
	return m.DeleteFn(ctx, id, requesterID)
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
			case 3:
				errNum = 0
				return "", domain.ErrForbidden
			}
			return ID_OK, nil
		},
	}
	handler := NewHealthRecordHandler(mockHealthRecordService)

	tests := []struct {
		test_name       string
		body            string
		expected_status int
		err_num         int
		requester_id    string
		requester_num   int8
	}{
		{
			test_name:       "Test OK",
			body:            make_body(USER_ID_REGULAR),
			expected_status: http.StatusCreated,
			err_num:         0,
			requester_id:    USER_ID_REGULAR,
			requester_num:   0,
		},
		{
			test_name:       "Test Invalid ID",
			body:            make_body(USER_ID_REGULAR),
			expected_status: http.StatusBadRequest,
			err_num:         1,
			requester_id:    USER_ID_REGULAR,
			requester_num:   0,
		},
		{
			test_name:       "Test Generic Error",
			body:            make_body(USER_ID_REGULAR),
			expected_status: http.StatusInternalServerError,
			err_num:         2,
			requester_id:    USER_ID_REGULAR,
			requester_num:   0,
		},
		{
			test_name:       "Test Forbidden Error",
			body:            make_body(USER_ID_REGULAR),
			expected_status: http.StatusInternalServerError,
			err_num:         2,
			requester_id:    USER_ID_REGULAR,
			requester_num:   0,
		},
		{
			test_name:       "Test Forbidden User ID",
			body:            make_body(USER_ID_MANAGER),
			expected_status: http.StatusForbidden,
			err_num:         0,
			requester_id:    USER_ID_REGULAR,
			requester_num:   0,
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
			requester_id:    USER_ID_REGULAR,
			requester_num:   0,
		},
		{
			test_name:       "Missing Requester ID",
			body:            make_body(USER_ID_REGULAR),
			expected_status: http.StatusForbidden,
			err_num:         2,
			requester_num:   0,
		},

		{
			test_name:       "Test Conversion Fail",
			body:            make_body(USER_ID_MANAGER),
			expected_status: http.StatusForbidden,
			err_num:         0,
			requester_num:   123,
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

			if tt.requester_id != "" {
				c.Set("userID", tt.requester_id)
			}

			if tt.requester_num == 123 {
				c.Set("userID", tt.requester_num)
			}

			errNum = tt.err_num
			handler.Create(c)
			status := c.Writer.Status()

			if status != tt.expected_status {
				t.Errorf("Unexpected status %v", status)
			}

		})
	}

}

func make_body(requesterID string) string {
	return fmt.Sprintf(`
			{
				"user_id": "%s",
				"health_record_type_id": "type",
				"value": 3.0,
				"notes": "Notes Ok",
				"recorded_at": "2026-06-25T19:30:00Z"
			}
			`, requesterID)
}

func TestHealthRecord_Update(t *testing.T) {
	mockHealthRecordService := &MockHealthRecordService{
		UpdateFn: func(ctx context.Context, id string, update_input servicedto.HealthRecordUpdateInput, userID string) error {
			if id == ID_NOT_FOUND {
				return domain.ErrHealthRecordNotFound
			}

			if id == ID_GENERIC_ERROR {
				return errors.New("Generic Error")
			}

			if id == ID_FORBIDDEN {
				return domain.ErrForbidden
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
		requester_id    string
		requester_num   int8
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
			requester_id:    USER_ID_REGULAR,
			requester_num:   0,
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
			requester_id:    USER_ID_REGULAR,
			requester_num:   0,
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
			requester_id:    USER_ID_REGULAR,
			requester_num:   0,
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
			requester_id:    USER_ID_REGULAR,
			requester_num:   0,
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
			requester_id:    USER_ID_REGULAR,
			requester_num:   0,
		},
		{
			test_name: "Test Generic Error",
			id:        ID_FORBIDDEN,
			body: `
			{
				"value": 3.0,
				"notes": "Notes Ok",
				"recorded_at": "2026-06-25T19:30:00Z"
			}
			`,
			expected_status: http.StatusForbidden,
			requester_id:    USER_ID_REGULAR,
			requester_num:   0,
		},
		{
			test_name: "Missing Requester ID",
			id:        ID_GENERIC_ERROR,
			body: `
			{
				"value": 3.0,
				"notes": "Notes Ok",
				"recorded_at": "2026-06-25T19:30:00Z"
			}
			`,
			expected_status: http.StatusForbidden,
			requester_num:   0,
		},
		{
			test_name: "Conversion Fail",
			id:        ID_GENERIC_ERROR,
			body: `
			{
				"value": 3.0,
				"notes": "Notes Ok",
				"recorded_at": "2026-06-25T19:30:00Z"
			}
			`,
			expected_status: http.StatusForbidden,
			requester_num:   123,
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

			if tt.requester_id != "" {
				c.Set("userID", USER_ID_REGULAR)
			}

			if tt.requester_num == 123 {
				c.Set("userID", tt.requester_num)
			}

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
		GetByIDFn: func(ctx context.Context, id string, userID string) (*domain.HealthRecord, error) {
			if id == ID_NOT_FOUND {
				return nil, domain.ErrHealthRecordNotFound
			}
			if id == ID_GENERIC_ERROR {
				return nil, errors.New("Generic Error")
			}
			if id == ID_NIL_RECORD {
				return nil, nil
			}
			return &domain.HealthRecord{
				UserID: USER_ID_REGULAR,
			}, nil
		},
	}

	handler := NewHealthRecordHandler(mockHealthRecordService)

	tests := []struct {
		test_name       string
		id              string
		expected_status int
		requester_id    string
		requester_num   int8
	}{
		{
			test_name:       "Test OK",
			id:              ID_OK,
			expected_status: http.StatusOK,
			requester_id:    USER_ID_REGULAR,
			requester_num:   0,
		},
		{
			test_name:       "Test Not Found",
			id:              ID_NOT_FOUND,
			expected_status: http.StatusNotFound,
			requester_id:    USER_ID_REGULAR,
			requester_num:   0,
		},
		{
			test_name:       "Test Generic Error",
			id:              ID_GENERIC_ERROR,
			expected_status: http.StatusInternalServerError,
			requester_id:    USER_ID_REGULAR,
			requester_num:   0,
		},
		{
			test_name:       "Test Generic Error",
			id:              ID_NIL_RECORD,
			expected_status: http.StatusInternalServerError,
			requester_id:    USER_ID_REGULAR,
			requester_num:   0,
		},
		{
			test_name:       "Test missing Id",
			expected_status: http.StatusBadRequest,
			requester_id:    USER_ID_REGULAR,
			requester_num:   0,
		},
		{
			test_name:       "Test Missing Requester ID",
			id:              ID_OK,
			expected_status: http.StatusForbidden,
			requester_num:   0,
		},
		{
			test_name:       "Test Conversion Fail",
			id:              ID_OK,
			expected_status: http.StatusForbidden,
			requester_num:   123,
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

			if tt.requester_id != "" {
				c.Set("userID", tt.requester_id)
			}

			if tt.requester_num == 123 {
				c.Set("userID", tt.requester_num)
			}
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
		ListByUserIDFn: func(ctx context.Context, userId string, requesterID string) ([]*domain.HealthRecord, error) {
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

			if userId == ID_FORBIDDEN {
				return nil, domain.ErrForbidden
			}
			if userId == ID_NIL_RECORD {
				var listRecord []*domain.HealthRecord
				listRecord = append(listRecord, nil)
				return listRecord, nil
			}
			return []*domain.HealthRecord{
				{
					ID:                 ID_OK,
					UserID:             USER_ID_REGULAR,
					HealthRecordTypeID: ID_OK,
					Value:              3.0,
					Notes:              &notes,
					RecordedAt:         time.Now().Add(-3 * time.Minute),
					CreatedAt:          time.Now().Add(-2 * time.Minute),
					UpdatedAt:          time.Now().Add(time.Minute),
				},
				{
					ID:                 ID_OK,
					UserID:             USER_ID_REGULAR,
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
		requester_id    string
		requester_num   int8
	}{
		{
			test_name:       "Test Ok",
			user_id:         ID_OK,
			expected_status: http.StatusOK,
			requester_id:    USER_ID_REGULAR,
			requester_num:   0,
		},
		{
			test_name:       "Test Nil Record",
			user_id:         ID_NIL_RECORD,
			expected_status: http.StatusOK,
			requester_id:    USER_ID_REGULAR,
			requester_num:   0,
		},
		{
			test_name:       "Test Not Found",
			user_id:         ID_NOT_FOUND,
			expected_status: http.StatusNotFound,
			requester_id:    USER_ID_REGULAR,
			requester_num:   0,
		},
		{
			test_name:       "Test Generic Error",
			user_id:         ID_GENERIC_ERROR,
			expected_status: http.StatusInternalServerError,
			requester_id:    USER_ID_REGULAR,
			requester_num:   0,
		},
		{
			test_name:       "Test Generic Error",
			user_id:         ID_EMPTY_RETURN,
			expected_status: http.StatusNotFound,
			requester_id:    USER_ID_REGULAR,
			requester_num:   0,
		},
		{
			test_name:       "Test Forbidden",
			user_id:         ID_FORBIDDEN,
			expected_status: http.StatusForbidden,
			requester_id:    USER_ID_REGULAR,
			requester_num:   0,
		},
		{
			test_name:       "Test Empty UserID",
			user_id:         "",
			expected_status: http.StatusBadRequest,
			requester_id:    USER_ID_REGULAR,
			requester_num:   0,
		},
		{
			test_name:       "Test Missing Requester ID",
			user_id:         ID_OK,
			expected_status: http.StatusForbidden,
			requester_num:   0,
		},
		{
			test_name:       "Test Missing Requester ID",
			user_id:         ID_OK,
			expected_status: http.StatusForbidden,
			requester_num:   123,
		},
	}

	for _, tt := range tests {
		t.Run(tt.test_name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			if tt.requester_id != "" {
				c.Set("userID", tt.requester_id)
			}

			if tt.requester_num == 123 {
				c.Set("userID", tt.requester_num)
			}

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
		DeleteFn: func(ctx context.Context, id string, requesterID string) error {
			if id == ID_NOT_FOUND {
				return domain.ErrHealthRecordNotFound
			}
			if id == ID_FORBIDDEN {
				return domain.ErrForbidden
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
		requester_id    string
		requester_num   int8
	}{
		{
			test_name:       "Test Ok",
			id:              ID_OK,
			expected_status: http.StatusNoContent,
			requester_id:    USER_ID_REGULAR,
			requester_num:   0,
		},
		{
			test_name:       "Test Not Found",
			id:              ID_NOT_FOUND,
			expected_status: http.StatusNotFound,
			requester_id:    USER_ID_REGULAR,
			requester_num:   0,
		},
		{
			test_name:       "Test Generic Error",
			id:              ID_GENERIC_ERROR,
			expected_status: http.StatusInternalServerError,
			requester_id:    USER_ID_REGULAR,
			requester_num:   0,
		},
		{
			test_name:       "Test Forbidden",
			id:              ID_FORBIDDEN,
			expected_status: http.StatusForbidden,
			requester_id:    USER_ID_REGULAR,
			requester_num:   0,
		},
		{
			test_name:       "Test No ID",
			id:              "",
			expected_status: http.StatusBadRequest,
			requester_id:    USER_ID_REGULAR,
			requester_num:   0,
		},
		{
			test_name:       "Test Missing Requester ID",
			id:              ID_OK,
			expected_status: http.StatusForbidden,
			requester_num:   0,
		},
		{
			test_name:       "Test conversion failed",
			id:              ID_OK,
			expected_status: http.StatusForbidden,
			requester_num:   123,
		},
	}
	for _, tt := range tests {
		t.Run(tt.test_name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			if tt.requester_id != "" {
				c.Set("userID", tt.requester_id)
			}

			if tt.requester_num == 123 {
				c.Set("userID", tt.requester_num)
			}

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
