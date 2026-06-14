package service

import (
	"context"
	"errors"
	"testing"

	"github.com/andrebarone77/cardiaflow-api/internal/domain"
	servicedto "github.com/andrebarone77/cardiaflow-api/internal/service/dto"
)

type MockUserRepository struct {
	SaveFn       func(ctx context.Context, user *domain.User) error
	GetByEmailFn func(ctx context.Context, email string) (*domain.User, error)
	GetByIdFn    func(ctx context.Context, id string) (*domain.User, error)
	DeleteFn     func(ctx context.Context, id string) error
}

func (m *MockUserRepository) Save(ctx context.Context, user *domain.User) error {
	return m.SaveFn(ctx, user)
}
func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return m.GetByEmailFn(ctx, email)
}

func (m *MockUserRepository) GetById(ctx context.Context, id string) (*domain.User, error) {
	return m.GetByIdFn(ctx, id)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	return m.DeleteFn(ctx, id)
}

func TestUserService_Create(t *testing.T) {

	mockRepo := &MockUserRepository{
		SaveFn: func(ctx context.Context, user *domain.User) error {
			if user.Email == "already_exists_user@email.com" {
				return domain.ErrEmailAlreadyExists
			}
			if user.PasswordHash == "" {
				return errors.New("password is empty")
			}
			return nil
		},
	}

	service := NewUserService(mockRepo)

	tests := []struct {
		test_name     string
		name          string
		email         string
		password      string
		saveError     error
		expects_error bool
	}{
		{
			test_name:     "Sucess",
			name:          "Andre",
			email:         "andre@email.com",
			password:      "password123",
			expects_error: false,
		},
		{
			test_name:     "User Exists",
			name:          "Andre",
			email:         "already_exists_user@email.com",
			password:      "password123",
			expects_error: true,
		},
		{
			test_name:     "Empty password",
			name:          "Andre",
			email:         "already_exists_user@email.com",
			expects_error: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.test_name, func(t *testing.T) {
			input := servicedto.CreateUserInput{
				Name:     tt.name,
				Email:    tt.email,
				Password: tt.password,
			}

			_, err := service.Create(context.Background(), input)
			if tt.expects_error && err == nil {
				t.Errorf("expected error, got nil")
			}

			if !tt.expects_error && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

		})
	}
}

func TestUserService_GetByEmail(t *testing.T) {
	mockRepo := &MockUserRepository{
		GetByEmailFn: func(ctx context.Context, email string) (*domain.User, error) {
			if email == "" {
				return nil, domain.ErrUserNotFound
			}

			return nil, nil
		},
	}

	service := NewUserService(mockRepo)

	tests := []struct {
		test_name     string
		email         string
		expects_error bool
	}{
		{
			test_name:     "success",
			email:         "andre@email.com",
			expects_error: false,
		},
		{
			test_name:     "missing_email",
			expects_error: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.test_name, func(t *testing.T) {

			_, err := service.GetByEmail(context.Background(), tt.email)

			if tt.expects_error && err == nil {
				t.Errorf("expected error, got nil")
			}

			if !tt.expects_error && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

		})
	}
}

func TestUserService_GetUserById(t *testing.T) {
	mockRepo := &MockUserRepository{
		GetByIdFn: func(ctx context.Context, id string) (*domain.User, error) {
			if id == UUID_NOT_FOUND {
				return nil, domain.ErrUserNotFound
			}
			return nil, nil

		},
	}

	service := NewUserService(mockRepo)

	tests := []struct {
		test_name     string
		id            string
		expects_error bool
	}{
		{
			test_name:     "success",
			id:            HEALTH_RECORD_NAME_OK,
			expects_error: false,
		},
		{
			test_name:     "not found",
			id:            UUID_NOT_FOUND,
			expects_error: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.test_name, func(t *testing.T) {
			_, err := service.GetById(context.Background(), tt.id)

			if !tt.expects_error && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

		})
	}

}

func TestUserService_Update(t *testing.T) {
	saveCalled := false
	getCalled := false
	mockRepo := &MockUserRepository{
		GetByIdFn: func(ctx context.Context, id string) (*domain.User, error) {
			getCalled = true
			if id == UUID_NOT_FOUND {
				return nil, domain.ErrUserNotFound
			}

			return &domain.User{
				ID:    id,
				Name:  "Old Name",
				Email: "old@email.com",
			}, nil

		},
		SaveFn: func(ctx context.Context, user *domain.User) error {
			saveCalled = true
			if user.ID == UUID_SAVE_ERRROR {
				return errors.New("Error saving")
			}
			return nil
		},
	}

	service := NewUserService(mockRepo)

	tests := []struct {
		test_name     string
		id            string
		name          string
		email         string
		password      string
		saveError     error
		expects_error bool
		expects_save  bool
		expects_get   bool
	}{
		{
			test_name:     "Sucess",
			id:            HEALTH_RECORD_NAME_OK,
			name:          "Andre",
			email:         "andre@email.com",
			password:      "password123",
			expects_error: false,
			expects_save:  true,
			expects_get:   true,
		},
		{
			test_name:     "Not Found",
			id:            UUID_NOT_FOUND,
			name:          "Andre",
			email:         "andre@email.com",
			password:      "password123",
			expects_error: true,
			expects_save:  false,
			expects_get:   true,
		},
		{
			test_name:     "Save Error",
			id:            UUID_SAVE_ERRROR,
			name:          "Andre",
			email:         "andre@email.com",
			password:      "password123",
			expects_error: true,
			expects_save:  true,
			expects_get:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.test_name, func(t *testing.T) {
			input := servicedto.UpdateUserInput{
				Name:     &tt.name,
				Email:    &tt.email,
				Password: &tt.password,
			}
			_, err := service.Update(context.Background(), tt.id, input)

			if !tt.expects_error && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.expects_error && err == nil {
				t.Errorf("expected error, got nil")
			}
			if tt.expects_save && !saveCalled {
				t.Errorf("Save not called")
			}
			if !tt.expects_save && saveCalled {
				t.Errorf("Save called when should not")
			}
			if tt.expects_get && !getCalled {
				t.Errorf("Get not called")
			}
			if !tt.expects_get && getCalled {
				t.Errorf("Get called when should not")
			}
			saveCalled = false
			getCalled = false
		})
	}

}

func TestUserService_Delete(t *testing.T) {
	deleteCalled := false
	mockRepo := &MockUserRepository{
		DeleteFn: func(ctx context.Context, id string) error {
			deleteCalled = true
			if id == UUID_NOT_FOUND {
				return domain.ErrUserNotFound
			}
			return nil
		},
	}
	service := NewUserService(mockRepo)

	tests := []struct {
		test_name             string
		id                    string
		expects_delete_called bool
		expect_error          bool
	}{
		{
			test_name:             "Sucess",
			id:                    HEALTH_RECORD_NAME_OK,
			expects_delete_called: true,
			expect_error:          false,
		},
		{
			test_name:             "Not Found",
			id:                    UUID_NOT_FOUND,
			expects_delete_called: true,
			expect_error:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.test_name, func(t *testing.T) {
			err := service.Delete(context.Background(), tt.id)

			if tt.expects_delete_called && !deleteCalled {
				t.Errorf("Delete function not called")
			}

			if !tt.expect_error && err != nil {
				t.Errorf("Unexpcted Error %v", err)
			}

			if tt.expect_error && err == nil {
				t.Errorf("Expecting error")
			}

		})
		deleteCalled = false
	}

}
