package service

import (
	"context"
	"testing"

	"github.com/andrebarone77/cardiaflow-api/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

const EMAIL_OK = "test_ok@email.com"
const EMAIL_NOT_FOUND = "test_not_found@email.com"

type MockAuthRepository struct {
	GetByEmailFn func(ctx context.Context, email string) (*domain.User, error)
}

func (m *MockAuthRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return m.GetByEmailFn(ctx, email)
}

func TestAuthService_Login(t *testing.T) {

	hash, _ := bcrypt.GenerateFromPassword(
		[]byte("ok"),
		bcrypt.DefaultCost,
	)

	mockRepo := &MockAuthRepository{
		GetByEmailFn: func(ctx context.Context, email string) (*domain.User, error) {
			if email == EMAIL_NOT_FOUND {
				return nil, domain.ErrUserNotFound
			}
			return &domain.User{
				ID:           UUID_OK,
				Email:        email,
				PasswordHash: string(hash),
			}, nil
		},
	}

	service := NewAuthService(mockRepo)

	tests := []struct {
		test_name  string
		email      string
		password   string
		expect_err bool
	}{
		{
			test_name:  "Test Ok",
			email:      EMAIL_OK,
			password:   "ok",
			expect_err: false,
		},
		{
			test_name:  "Test OkNot Found",
			email:      EMAIL_NOT_FOUND,
			password:   "ok",
			expect_err: true,
		},
		{
			test_name:  "Test Hash not Ok",
			email:      EMAIL_OK,
			password:   "notok",
			expect_err: true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.test_name, func(t *testing.T) {
			t.Setenv("JWT_SECRET", "my-test-secret")
			t.Setenv("JWT_EXPIRES_IN", "24h")

			token, err := service.Login(context.Background(), tt.email, tt.password)

			if !tt.expect_err && err != nil {
				t.Errorf("Unexpected Error %v", err)
			}

			if !tt.expect_err && (token == "") {
				t.Errorf("Return is empty")
			}

		})

	}
}
