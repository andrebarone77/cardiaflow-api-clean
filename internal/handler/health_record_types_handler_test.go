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

type MockHealthRecordType struct {
	CreateFn    func(ctx context.Context, input servicedto.HealthRecordTypeInput) (string, error)
	GetByIDFn   func(ctx context.Context, id string) (*domain.HealthRecordType, error)
	GetByCodeFn func(ctx context.Context, code string) (*domain.HealthRecordType, error)
	GetAllFn    func(ctx context.Context) ([]*domain.HealthRecordType, error)
	DeleteFn    func(ctx context.Context, id string) error
	UpdateFn    func(ctx context.Context, id string, update servicedto.HealthRecordTypeUpdateInput) (*domain.HealthRecordType, error)
}

func (m *MockHealthRecordType) Create(ctx context.Context, input servicedto.HealthRecordTypeInput) (string, error) {
	return m.CreateFn(ctx, input)
}

func (m *MockHealthRecordType) GetByID(ctx context.Context, id string) (*domain.HealthRecordType, error) {
	return m.GetByIDFn(ctx, id)
}

func (m *MockHealthRecordType) GetByCode(ctx context.Context, code string) (*domain.HealthRecordType, error) {
	return m.GetByCodeFn(ctx, code)
}

func (m *MockHealthRecordType) GetAll(ctx context.Context) ([]*domain.HealthRecordType, error) {
	return m.GetAllFn(ctx)
}
func (m *MockHealthRecordType) Delete(ctx context.Context, id string) error {
	return m.DeleteFn(ctx, id)
}

func (m *MockHealthRecordType) Update(ctx context.Context, id string, update servicedto.HealthRecordTypeUpdateInput) (*domain.HealthRecordType, error) {
	return m.UpdateFn(ctx, id, update)
}

func TestHealhRecordTypeHandler_Create(t *testing.T) {
	errNum := 0
	mockHealhRecordTypeService := &MockHealthRecordType{
		CreateFn: func(ctx context.Context, input servicedto.HealthRecordTypeInput) (string, error) {
			switch errNum {
			case 1:
				errNum = 0
				return "", domain.ErrCodeRequired
			case 2:
				errNum = 0
				return "", domain.ErrCodeTooLong
			case 3:
				errNum = 0
				return "", domain.ErrHealthRecordTypeAlreadyExists
			case 4:
				errNum = 0
				return "", domain.ErrCodeInvalid
			case 5:
				errNum = 0
				return "", errors.New("Generic error")
			}

			return "", nil
		},
	}

	handler := NewHealthRecordTypeHandler(mockHealhRecordTypeService)

	tests := []struct {
		test_name       string
		body            string
		err_num         int
		expected_status int
	}{
		{
			test_name: "Test OK",
			body: `
				{
					"name":"teste",
					"code":"code",
					"unit":"unit"	
				}
			`,
			err_num:         0,
			expected_status: http.StatusCreated,
		},
		{
			test_name: "Test Malformed Body",
			body: `
				{
					"name":"teste,
					"code":"code",
					"unit":"unit"	
				}
			`,
			err_num:         0,
			expected_status: http.StatusBadRequest,
		},
		{
			test_name: "Test Code Required",
			body: `
				{
					"name":"teste",
					"code":"code",
					"unit":"unit"	
				}
			`,
			err_num:         1,
			expected_status: http.StatusBadRequest,
		},
		{
			test_name: "Test Code Already Exists",
			body: `
				{
					"name":"teste",
					"code":"code",
					"unit":"unit"	
				}
			`,
			err_num:         3,
			expected_status: http.StatusConflict,
		},
		{
			test_name: "Test Code Too Long",
			body: `
				{
					"name":"teste",
					"code":"codecodecodecodecodecodecodecodecodecodecodecodecodecode",
					"unit":"unit"	
				}
			`,
			err_num:         2,
			expected_status: http.StatusBadRequest,
		},
		{
			test_name: "Test OK",
			body: `
				{
					"name":"teste",
					"code":"Code",
					"unit":"unit"	
				}
			`,
			err_num:         4,
			expected_status: http.StatusBadRequest,
		},
		{
			test_name: "Test OK",
			body: `
				{
					"name":"teste",
					"code":"Code",
					"unit":"unit"	
				}
			`,
			err_num:         5,
			expected_status: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.test_name, func(t *testing.T) {
			errNum = tt.err_num
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(
				http.MethodPost,
				"/health_record_type",
				strings.NewReader(tt.body),
			)
			handler.Create(c)
			status := c.Writer.Status()

			if status != tt.expected_status {
				t.Errorf("Expected Status %v received %v", tt.expected_status, status)
			}
		})
	}
}

func TestHealhRecordTypeHandler_GetById(t *testing.T) {
	errNum := 0
	mockHealhRecordTypeService := &MockHealthRecordType{
		GetByIDFn: func(ctx context.Context, id string) (*domain.HealthRecordType, error) {
			switch errNum {
			case 1:
				errNum = 0
				return nil, domain.ErrHealthRecordTypeNotFound
			case 2:
				errNum = 0
				return nil, errors.New("Generic Error")
			}
			unit := "mm"
			return &domain.HealthRecordType{
				ID:        ID_OK,
				Name:      "code",
				Code:      "code",
				Unit:      &unit,
				IsSystem:  false,
				CreatedAt: time.Now().Add(-5 * time.Minute),
				UpdatedAt: time.Now(),
			}, nil
		},
	}

	handler := NewHealthRecordTypeHandler(mockHealhRecordTypeService)

	tests := []struct {
		test_name       string
		id              string
		err_num         int
		expected_status int
	}{
		{
			test_name:       "Test OK",
			id:              ID_OK,
			err_num:         0,
			expected_status: http.StatusOK,
		},
		{
			test_name:       "Test Not Found",
			id:              ID_NOT_FOUND,
			err_num:         1,
			expected_status: http.StatusNotFound,
		},
		{
			test_name:       "Test OK",
			id:              ID_GENERIC_ERROR,
			err_num:         2,
			expected_status: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.test_name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			errNum = tt.err_num

			c.Request = httptest.NewRequest(
				http.MethodGet,
				fmt.Sprintf("/health_record_type/%v", tt.id),
				nil,
			)
			c.Params = gin.Params{
				{
					Key:   "id",
					Value: tt.id,
				},
			}
			handler.GetByID(c)
			status := c.Writer.Status()

			if status != tt.expected_status {
				t.Errorf("Expected Status %v received %v", tt.expected_status, status)
			}

		})
	}

}

func TestHealthRecordType_GetByCode(t *testing.T) {
	mmUnit := "mm"
	errNum := 0
	mockHealthRecordType := &MockHealthRecordType{
		GetByCodeFn: func(ctx context.Context, code string) (*domain.HealthRecordType, error) {

			switch errNum {
			case 1:
				errNum = 0
				return nil, domain.ErrHealthRecordTypeNotFound
			case 2:
				errNum = 0
				return nil, errors.New("Generic Error")

			}

			return &domain.HealthRecordType{
				ID:        ID_OK,
				Name:      "Code",
				Code:      "code",
				Unit:      &mmUnit,
				CreatedAt: time.Now().Add(-5 * time.Minute),
				UpdatedAt: time.Now(),
				IsSystem:  false,
			}, nil
		},
	}

	handler := NewHealthRecordTypeHandler(mockHealthRecordType)

	tests := []struct {
		test_name       string
		code            string
		err_num         int
		expected_status int
	}{
		{
			test_name:       "Test OK",
			code:            "code",
			err_num:         0,
			expected_status: http.StatusOK,
		},
		{
			test_name:       "Test Not Found",
			code:            "code",
			err_num:         1,
			expected_status: http.StatusNotFound,
		},
		{
			test_name:       "Test GenericError",
			code:            "code",
			err_num:         2,
			expected_status: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/api/healthrecordtype/code/%v", tt.code),
			nil,
		)

		errNum = tt.err_num

		c.Params = gin.Params{
			{
				Key:   "code",
				Value: tt.code,
			},
		}

		handler.GetByCode(c)
		status := c.Writer.Status()

		if tt.expected_status != status {
			t.Errorf("Unexpected Error: %v", status)
		}
	}

}

func TestHealhRecordTypeHandler_GetAll(t *testing.T) {
	errNum := 0
	mmUnit := "mm"
	MockHealthRecordType := &MockHealthRecordType{
		GetAllFn: func(ctx context.Context) ([]*domain.HealthRecordType, error) {
			switch errNum {
			case 1:
				errNum = 0
				return nil, domain.ErrHealthRecordTypeNotFound
			case 2:
				errNum = 0
				return nil, errors.New("Generic Error")

			}
			record1 := &domain.HealthRecordType{
				ID:        ID_OK,
				Name:      "Code",
				Code:      "code",
				Unit:      &mmUnit,
				CreatedAt: time.Now().Add(-5 * time.Minute),
				UpdatedAt: time.Now(),
				IsSystem:  false,
			}

			var hrt []*domain.HealthRecordType
			hrt = append(hrt, record1)
			return hrt, nil
		},
	}

	handler := NewHealthRecordTypeHandler(MockHealthRecordType)

	tests := []struct {
		test_name       string
		err_num         int
		expected_status int
	}{
		{
			test_name:       "Test Ok",
			err_num:         0,
			expected_status: http.StatusOK,
		},
		{
			test_name:       "Test Not Found",
			err_num:         1,
			expected_status: http.StatusNotFound,
		},
		{
			test_name:       "Test Generic Error",
			err_num:         2,
			expected_status: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(
			http.MethodGet,
			"/api/healthrecordtypes",
			nil,
		)
		errNum = tt.err_num
		handler.GetAll(c)
		status := c.Writer.Status()

		if status != tt.expected_status {
			t.Errorf("Unexpected status %v", status)
		}

	}
}

func TestHealthHandlerType_Update(t *testing.T) {
	mmUnit := "mm"
	errNum := 0
	mockHealthHandlerType := &MockHealthRecordType{
		UpdateFn: func(ctx context.Context, id string, update servicedto.HealthRecordTypeUpdateInput) (*domain.HealthRecordType, error) {
			switch errNum {
			case 1:
				errNum = 0
				return nil, domain.ErrHealthRecordTypeNotFound
			case 2:
				errNum = 0
				return nil, domain.ErrHealthRecordTypeAlreadyExists
			case 3:
				errNum = 0
				return nil, domain.ErrHealthRecordTypeImmutable
			case 4:
				errNum = 0
				return nil, errors.New("Generic Error")
			}
			return &domain.HealthRecordType{
				ID:        ID_OK,
				Name:      "Code",
				Code:      "code",
				Unit:      &mmUnit,
				CreatedAt: time.Now().Add(-5 * time.Minute),
				UpdatedAt: time.Now(),
				IsSystem:  false,
			}, nil
		},
	}

	handler := NewHealthRecordTypeHandler(mockHealthHandlerType)

	tests := []struct {
		test_name       string
		id              string
		body            string
		err_num         int
		expected_status int
	}{
		{
			test_name: "Test OK",
			id:        ID_OK,
			body: `
			{
				"name":"name",
				"code":"code",
				"unit":"unit"
			}
			`,
			err_num:         0,
			expected_status: http.StatusOK,
		},
		{
			test_name: "Test Not Found",
			id:        ID_OK,
			body: `
			{
				"name":"name",
				"code":"code",
				"unit":"unit"
			}
			`,
			err_num:         1,
			expected_status: http.StatusNotFound,
		},
		{
			test_name: "Test HealthRecordType Already Exists",
			id:        ID_OK,
			body: `
			{
				"name":"name",
				"code":"code",
				"unit":"unit"
			}
			`,
			err_num:         2,
			expected_status: http.StatusConflict,
		},

		{
			test_name: "Test HealthRecordType Immutable",
			id:        ID_OK,
			body: `
			{
				"name":"name",
				"code":"code",
				"unit":"unit"
			}
			`,
			err_num:         3,
			expected_status: http.StatusForbidden,
		},

		{
			test_name: "Test HealthRecordType Generic Error",
			id:        ID_OK,
			body: `
			{
				"name":"name",
				"code":"code",
				"unit":"unit"
			}
			`,
			err_num:         4,
			expected_status: http.StatusInternalServerError,
		},
		{
			test_name: "Test No Id",
			id:        "",
			body: `
			{
				"name":"name",
				"code":"code",
				"unit":"unit"
			}
			`,
			err_num:         0,
			expected_status: http.StatusBadRequest,
		},
		{
			test_name: "Test No Body",
			id:        ID_OK,
			body: `
			{

			}
			`,
			err_num:         0,
			expected_status: http.StatusBadRequest,
		},
		{
			test_name: "Test Malformed Body",
			id:        ID_OK,
			body: `
				{
					"name":"teste,
					"code":"code",
					"unit":"unit"	
				}
			`,
			err_num:         0,
			expected_status: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.test_name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(
				http.MethodPatch,
				fmt.Sprintf("/api/healthrecordtypes/%v", tt.id),
				strings.NewReader(tt.body),
			)
			c.Params = gin.Params{
				{Key: "id",
					Value: tt.id},
			}
			errNum = tt.err_num
			handler.Update(c)
			status := c.Writer.Status()

			if status != tt.expected_status {
				t.Errorf("Unexpected Status %v", status)
			}

		})
	}

}
func TestHealhRecordTypeHandler_Delete(t *testing.T) {
	mockHealthRecordType := &MockHealthRecordType{
		DeleteFn: func(ctx context.Context, id string) error {
			if id == ID_NOT_FOUND {
				return domain.ErrHealthRecordTypeNotFound
			}
			if id == ID_IMMUTABLE {
				return domain.ErrHealthRecordTypeImmutable
			}
			return nil
		},
	}

	handler := NewHealthRecordTypeHandler(mockHealthRecordType)

	tests := []struct {
		test_name       string
		id              string
		expected_status int
	}{
		{
			test_name:       "Test OK",
			id:              ID_OK,
			expected_status: http.StatusNoContent,
		},
		{
			test_name:       "Test NotFound",
			id:              ID_NOT_FOUND,
			expected_status: http.StatusNotFound,
		},
		{
			test_name:       "Test Immutable",
			id:              ID_IMMUTABLE,
			expected_status: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.test_name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(
				http.MethodDelete,
				fmt.Sprintf("/api/healthrecordtypes?id=%v", tt.id),
				nil,
			)

			handler.Delete(c)
			status := c.Writer.Status()

			if status != tt.expected_status {
				t.Errorf("Unexpcetd Status %v", status)
			}

		})
	}

}
