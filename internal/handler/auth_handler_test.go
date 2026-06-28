package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/andrebarone77/cardiaflow-api/internal/domain"
	"github.com/gin-gonic/gin"
)

type MockAuthService struct {
	LoginFn func(ctx context.Context, email string, password string) (string, error)
}

func (m *MockAuthService) Login(ctx context.Context, email string, password string) (string, error) {
	return m.LoginFn(ctx, email, password)
}

func TestAuthHandler_Login(t *testing.T) {
	mockAuthService := &MockAuthService{
		LoginFn: func(ctx context.Context, email, password string) (string, error) {
			if email == EMAIL_NOT_AUTHORIZED {
				return "", domain.ErrNotAuthorized
			}
			return "", nil
		},
	}

	handler := NewAuthHandler(mockAuthService)

	tests := []struct {
		test_name string
		body      string

		expected_status int
	}{
		{
			test_name: "Test Ok",
			body: `{
				"email":           "ok@test.com",
				"password":        "ok"
				}			`,
			expected_status: http.StatusOK,
		},
		{
			test_name: "Test Malformed Body",
			body: `{
				"email":           "malformed@test.com,
				"password":        "ok"
				}			`,
			expected_status: http.StatusBadRequest,
		},
		{
			test_name: "Test Malformed Body",
			body: `{
				"email":           "email@notauthorized.com",
				"password":        "ok"
				}			`,
			expected_status: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.test_name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(
				http.MethodPost,
				"/api/auth/login",
				strings.NewReader(tt.body),
			)
			c.Request.Header.Set(
				"Content-Type",
				"application/json",
			)

			handler.Login(c)
			status := c.Writer.Status()

			if status != tt.expected_status {
				t.Errorf("Unexpected Status %v", status)
			}

		})
	}
}
