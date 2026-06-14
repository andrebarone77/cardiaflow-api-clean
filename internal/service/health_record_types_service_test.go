package service

import (
	"context"
	"errors"
	"testing"

	"github.com/andrebarone77/cardiaflow-api/internal/domain"
	"github.com/andrebarone77/cardiaflow-api/internal/service/dto"
)

const HEALTH_RECORD_NAME_OK = "name_ok"
const HEALTH_RECORD_NAME_EXISTS = "name_exists"

type MockHealthRecordTypeRepository struct {
	CreateFn    func(ctx context.Context, h *domain.HealthRecordType) (string, error)
	GetByIdFn   func(ctx context.Context, id string) (*domain.HealthRecordType, error)
	GetByCodeFn func(ctx context.Context, code string) (*domain.HealthRecordType, error)
	GetAllFn    func(ctx context.Context) ([]*domain.HealthRecordType, error)
	UpdateFn    func(ctx context.Context, id string, h *domain.HealthRecordType) error
	DeleteFn    func(ctx context.Context, id string) error
	IsSystemFn  func(ctx context.Context, id string) (bool, error)
}

func (m *MockHealthRecordTypeRepository) Create(ctx context.Context, h *domain.HealthRecordType) (string, error) {
	return m.CreateFn(ctx, h)
}

func (m *MockHealthRecordTypeRepository) GetById(ctx context.Context, id string) (*domain.HealthRecordType, error) {
	return m.GetByIdFn(ctx, id)
}

func (m *MockHealthRecordTypeRepository) GetByCode(ctx context.Context, code string) (*domain.HealthRecordType, error) {
	return m.GetByCodeFn(ctx, code)
}

func (m *MockHealthRecordTypeRepository) GetAll(ctx context.Context) ([]*domain.HealthRecordType, error) {
	return m.GetAllFn(ctx)
}

func (m *MockHealthRecordTypeRepository) Update(ctx context.Context, id string, h *domain.HealthRecordType) error {
	return m.UpdateFn(ctx, id, h)
}

func (m *MockHealthRecordTypeRepository) Delete(ctx context.Context, id string) error {
	return m.DeleteFn(ctx, id)
}

func (m *MockHealthRecordTypeRepository) IsSystem(ctx context.Context, id string) (bool, error) {
	return m.IsSystemFn(ctx, id)
}

func TestHealthRecordType_Create(t *testing.T) {
	createCalled := false
	unit_ok := "unitok"
	mockRepo := &MockHealthRecordTypeRepository{
		CreateFn: func(ctx context.Context, h *domain.HealthRecordType) (string, error) {
			createCalled = true
			if h.Name == HEALTH_RECORD_NAME_EXISTS {
				return "", domain.ErrHealthRecordAlreadyExists
			}
			return "", nil
		},
	}

	service := NewHealthRecordTypeService(mockRepo)

	tests := []struct {
		test_name             string
		health_record         dto.HealthRecordTypeInput
		expect_error          bool
		expcted_create_called bool
	}{
		{
			test_name: "Success",
			health_record: dto.HealthRecordTypeInput{
				Name: HEALTH_RECORD_NAME_OK,
				Code: CODE_OK,
				Unit: &unit_ok,
			},
			expect_error:          false,
			expcted_create_called: true,
		},
		{
			test_name: "Already Exists",
			health_record: dto.HealthRecordTypeInput{
				Name: HEALTH_RECORD_NAME_EXISTS,
				Code: CODE_OK,
				Unit: &unit_ok,
			},
			expect_error:          true,
			expcted_create_called: true,
		},
		{
			test_name: "Code Missing",
			health_record: dto.HealthRecordTypeInput{
				Name: HEALTH_RECORD_NAME_OK,
				Unit: &unit_ok,
			},
			expect_error:          true,
			expcted_create_called: false,
		},
		{
			test_name: "Long Code Name",
			health_record: dto.HealthRecordTypeInput{
				Name: HEALTH_RECORD_NAME_OK,
				Code: LONG_CODE,
				Unit: &unit_ok,
			},
			expect_error:          true,
			expcted_create_called: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.test_name, func(t *testing.T) {
			_, err := service.Create(context.Background(), tt.health_record)

			if tt.expcted_create_called && !createCalled {
				t.Errorf("Error Create Not Called")
			}

			if !tt.expect_error && err != nil {
				t.Errorf("Unexpected Error: %v", err)
			}

			if tt.expect_error && err == nil {
				t.Errorf("Error expected")
			}

		})
		createCalled = false
	}

}

func TestHealthRecordType_GetById(t *testing.T) {
	getByIdCalled := false
	mockRepo := &MockHealthRecordTypeRepository{
		GetByIdFn: func(ctx context.Context, id string) (*domain.HealthRecordType, error) {
			getByIdCalled = true
			if id == UUID_NOT_FOUND {
				return nil, domain.ErrCodeInvalid
			}
			return nil, nil
		},
	}

	service := NewHealthRecordTypeService(mockRepo)

	tests := []struct {
		test_name        string
		id               string
		error_expect     bool
		get_by_id_called bool
	}{
		{
			test_name:        "success",
			id:               UUID_OK,
			error_expect:     false,
			get_by_id_called: true,
		},
		{
			test_name:        "not_found",
			id:               UUID_NOT_FOUND,
			error_expect:     true,
			get_by_id_called: true,
		},
	}

	for _, tt := range tests {
		_, err := service.GetByID(context.Background(), tt.id)
		t.Run(tt.test_name, func(t *testing.T) {

			if tt.get_by_id_called && !getByIdCalled {
				t.Errorf("GetById not called")
			}

			if !tt.error_expect && err != nil {
				t.Errorf("Error found %v", err)
			}

			if tt.error_expect && err == nil {
				t.Errorf("Error expected")
			}

		})
	}

}

func TestHealthRecordType_GetByCode(t *testing.T) {
	getByCodeCalled := false
	mockRepo := &MockHealthRecordTypeRepository{
		GetByCodeFn: func(ctx context.Context, code string) (*domain.HealthRecordType, error) {
			getByCodeCalled = true
			if code == INVALID_CODE {
				return nil, domain.ErrCodeInvalid
			}
			return nil, nil
		},
	}

	service := NewHealthRecordTypeService(mockRepo)

	tests := []struct {
		test_name        string
		code             string
		error_expect     bool
		get_by_id_called bool
	}{
		{
			test_name:        "success",
			code:             CODE_OK,
			error_expect:     false,
			get_by_id_called: true,
		},
		{
			test_name:        "not_found",
			code:             INVALID_CODE,
			error_expect:     true,
			get_by_id_called: true,
		},
	}

	for _, tt := range tests {
		_, err := service.GetByCode(context.Background(), tt.code)
		t.Run(tt.test_name, func(t *testing.T) {

			if tt.get_by_id_called && !getByCodeCalled {
				t.Errorf("GetById not called")
			}

			if !tt.error_expect && err != nil {
				t.Errorf("Error found %v", err)
			}

			if tt.error_expect && err == nil {
				t.Errorf("Error expected")
			}

		})
	}
}

func TestHealthRecordType_GetAll(t *testing.T) {
	getAllCalled := false
	counter := 0

	mockRepo := &MockHealthRecordTypeRepository{
		GetAllFn: func(ctx context.Context) ([]*domain.HealthRecordType, error) {
			getAllCalled = true
			if counter > 0 {
				return nil, errors.New("Failed to access DB")
			}
			counter++
			return nil, nil
		},
	}

	service := NewHealthRecordTypeService(mockRepo)

	tests := []struct {
		test_name      string
		get_all_called bool
		expect_error   bool
	}{
		{
			test_name:      "Sucess",
			get_all_called: true,
			expect_error:   false,
		}, {
			test_name:      "Fail",
			get_all_called: true,
			expect_error:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.test_name, func(t *testing.T) {
			_, err := service.GetAll(context.Background())

			if tt.get_all_called && !getAllCalled {
				t.Errorf("GetAll expected and not Called")
			}

			if !tt.expect_error && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if tt.expect_error && err == nil {
				t.Errorf("Error Expected")
			}
		})
		getAllCalled = false
	}

}

func TestHealthRecordType_Update(t *testing.T) {
	updateCalled := false
	isSystemCalled := false
	getByIdCalled := false

	name_ok := "name_ok"
	code_ok := "code_ok"
	code_is_system := UUID_IS_SYSTEM
	unit_ok := "unit_ok"

	healthRecordTypeIsSystem := domain.HealthRecordType{
		ID:       UUID_OK,
		Name:     "Generic HealthRecordType",
		Code:     code_ok,
		Unit:     &unit_ok,
		IsSystem: true,
	}
	healthRecordTypeNotSystem := domain.HealthRecordType{
		ID:       UUID_OK,
		Name:     "Generic HealthRecordType",
		Code:     code_ok,
		Unit:     &unit_ok,
		IsSystem: false,
	}

	mockRepo := &MockHealthRecordTypeRepository{
		UpdateFn: func(ctx context.Context, id string, h *domain.HealthRecordType) error {
			updateCalled = true
			if id == UUID_NOT_FOUND {
				return domain.ErrHealthRecordTypeNotFound
			}

			return nil
		},
		IsSystemFn: func(ctx context.Context, id string) (bool, error) {
			isSystemCalled = true

			if id == UUID_GENERIC_ERROR {
				return false, errors.New("Error on isSystem")
			}

			if id == UUID_IS_SYSTEM {
				return true, nil
			}

			return false, nil
		},
		GetByIdFn: func(ctx context.Context, id string) (*domain.HealthRecordType, error) {
			getByIdCalled = true
			if id == "" {
				return nil, domain.ErrCodeInvalid
			}
			if id == UUID_NOT_FOUND {
				return nil, domain.ErrHealthRecordTypeNotFound
			}
			if id == UUID_IS_SYSTEM {
				return &healthRecordTypeIsSystem, nil
			}
			return &healthRecordTypeNotSystem, nil
		},
	}

	service := NewHealthRecordTypeService(mockRepo)

	tests := []struct {
		test_name               string
		id                      string
		health_record_type      dto.HealthRecordTypeUpdateInput
		expect_error            bool
		expect_update_called    bool
		expect_is_system_called bool
		expect_get_by_id_called bool
	}{
		{
			test_name: "Sucess",
			id:        UUID_OK,
			health_record_type: dto.HealthRecordTypeUpdateInput{
				Name: &name_ok,
				Code: &code_ok,
				Unit: &unit_ok,
			},
			expect_error:            false,
			expect_update_called:    true,
			expect_is_system_called: true,
			expect_get_by_id_called: true,
		},
		{
			test_name: "Missing ID",
			health_record_type: dto.HealthRecordTypeUpdateInput{
				Name: &name_ok,
				Code: &code_ok,
				Unit: &unit_ok,
			},
			expect_error:            true,
			expect_update_called:    false,
			expect_is_system_called: true,
			expect_get_by_id_called: true,
		},
		{
			test_name: "ID Not Found",
			id:        UUID_NOT_FOUND,
			health_record_type: dto.HealthRecordTypeUpdateInput{
				Name: &name_ok,
				Code: &code_ok,
				Unit: &unit_ok,
			},
			expect_error:            true,
			expect_update_called:    true,
			expect_is_system_called: true,
			expect_get_by_id_called: true,
		},
		{
			test_name: "IsSystem",
			id:        UUID_IS_SYSTEM,
			health_record_type: dto.HealthRecordTypeUpdateInput{
				Name: &name_ok,
				Code: &code_is_system,
				Unit: &unit_ok,
			},
			expect_error:            true,
			expect_update_called:    false,
			expect_is_system_called: true,
			expect_get_by_id_called: false,
		},
		{
			test_name:               "Missing All Attributes",
			id:                      UUID_OK,
			health_record_type:      dto.HealthRecordTypeUpdateInput{},
			expect_error:            true,
			expect_update_called:    false,
			expect_is_system_called: true,
			expect_get_by_id_called: true,
		},
		{
			test_name:               "Missing All Attributes",
			id:                      UUID_GENERIC_ERROR,
			health_record_type:      dto.HealthRecordTypeUpdateInput{},
			expect_error:            true,
			expect_update_called:    false,
			expect_is_system_called: true,
			expect_get_by_id_called: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.test_name, func(t *testing.T) {
			_, err := service.Update(context.Background(), tt.id, tt.health_record_type)

			if tt.expect_is_system_called && !isSystemCalled {
				t.Errorf("IsSystem not called")
			}

			if tt.expect_get_by_id_called && !getByIdCalled {
				t.Errorf("GetByIdCalled not called")
			}
			if !tt.expect_update_called && updateCalled {
				t.Errorf("Expect Update called")
			}

			if !tt.expect_error && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if tt.expect_error && err == nil {
				t.Errorf("Expect Error")
			}

		})

		updateCalled = false
		isSystemCalled = false
		getByIdCalled = false
	}

}

func TestHealthRecordType_Delete(t *testing.T) {
	deleteCalled := false
	isSystemCalled := false
	mockRepo := &MockHealthRecordTypeRepository{
		DeleteFn: func(ctx context.Context, id string) error {
			deleteCalled = true
			if id == UUID_NOT_FOUND {
				return domain.ErrHealthRecordTypeNotFound
			}
			return nil
		}, IsSystemFn: func(ctx context.Context, id string) (bool, error) {
			isSystemCalled = true

			if id == UUID_IS_SYSTEM {
				return true, nil
			}

			if id == UUID_GENERIC_ERROR {
				return false, errors.New("Error on isSystem")
			}

			return false, nil
		},
	}

	service := NewHealthRecordTypeService(mockRepo)

	tests := []struct {
		test_name               string
		id                      string
		expect_error            bool
		expext_delete_called    bool
		expect_is_system_called bool
	}{
		{
			test_name:               "Success",
			id:                      UUID_OK,
			expect_error:            false,
			expext_delete_called:    true,
			expect_is_system_called: true,
		}, {
			test_name:               "Not Found",
			id:                      UUID_NOT_FOUND,
			expect_error:            true,
			expext_delete_called:    true,
			expect_is_system_called: true,
		},
		{
			test_name:               "Not Found",
			id:                      UUID_IS_SYSTEM,
			expect_error:            true,
			expext_delete_called:    false,
			expect_is_system_called: false,
		},
		{
			test_name:               "Not Found",
			id:                      UUID_GENERIC_ERROR,
			expect_error:            true,
			expext_delete_called:    false,
			expect_is_system_called: false,
		},
	}

	for _, tt := range tests {

		t.Run(tt.test_name, func(t *testing.T) {
			err := service.Delete(context.Background(), tt.id)

			if tt.expext_delete_called && !deleteCalled {
				t.Error("Expect Delete() to be called")
			}
			if tt.expect_is_system_called && !isSystemCalled {
				t.Error("Expect Delete() to be called")
			}
			if !tt.expect_error && err != nil {
				t.Errorf("Unexpect error: %v", err)
			}

			if tt.expect_error && err == nil {
				t.Errorf("Unexpect error")
			}

		})
		deleteCalled = false
		isSystemCalled = false
	}
}
