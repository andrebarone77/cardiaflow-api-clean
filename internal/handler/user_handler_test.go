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

const NAME_OK = "Teste name ok"
const EMAIL_OK = "email@ok.com"
const ID_OK = "eed53c96-3de2-40a2-afa7-f945f228e59c"
const PASSORD_OK = "Pa55w0rd"

const ID_GENERIC_ERROR = "e238e28a-c655-4c6f-a47a-2057004bac84"

const EMAIL_EXISTS = "email@exists.com"
const ID_EMAIL_EXISTS = "964951c7-93a9-4e34-9043-2876e6f8c148"

const EMAIL_OTHER_ERROR = "email@other.com"

const EMAIL_NOT_FOUND = "email@notfound.com"
const ID_NOT_FOUND = "3977f63b-86fe-44b8-8ff9-83261b212661"

type MockUserService struct {
	CreateFn     func(ctx context.Context, req servicedto.CreateUserInput) (*domain.User, error)
	GetByEmailFn func(ctx context.Context, email string) (*domain.User, error)
	GetByIdFn    func(ctx context.Context, id string) (*domain.User, error)
	DeleteFn     func(ctx context.Context, id string) error
	UpdateFn     func(ctx context.Context, id string, req servicedto.UpdateUserInput) (*domain.User, error)
}

func (m *MockUserService) Create(ctx context.Context, req servicedto.CreateUserInput) (*domain.User, error) {
	return m.CreateFn(ctx, req)
}

func (m *MockUserService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return m.GetByEmailFn(ctx, email)
}

func (m *MockUserService) GetById(ctx context.Context, id string) (*domain.User, error) {
	return m.GetByIdFn(ctx, id)
}

func (m *MockUserService) Delete(ctx context.Context, id string) error {
	return m.DeleteFn(ctx, id)
}

func (m *MockUserService) Update(ctx context.Context, id string, req servicedto.UpdateUserInput) (*domain.User, error) {
	return m.UpdateFn(ctx, id, req)
}

func TestUserHandler_Create(t *testing.T) {
	mockService := &MockUserService{
		CreateFn: func(ctx context.Context, req servicedto.CreateUserInput) (*domain.User, error) {

			if req.Email == EMAIL_EXISTS {
				return nil, domain.ErrEmailAlreadyExists
			}

			if req.Email == EMAIL_OTHER_ERROR {
				return nil, errors.New("Other error")
			}

			return &domain.User{
				ID:        ID_OK,
				Name:      NAME_OK,
				Email:     EMAIL_OK,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now().Add(3 * time.Minute),
			}, nil
		},
	}

	handler := NewUserHandler(mockService)

	tests := []struct {
		test_name    string
		body         string
		expect_error bool
	}{
		{
			test_name: "Test ok",
			body: `{
				"name":"Teste Ok",
				"email":"email@ok.com",
				"password":"Pa55w0rd"
				}
			`,
			expect_error: false,
		},
		{
			test_name: "Test missing name",
			body: `{
				"email":"email@ok.com",
				"password":"Pa55w0rd"
				}
			`,
			expect_error: true,
		},
		{
			test_name: "Test email exists",
			body: `{
				"name":"Email exists",
				"email":"email@exists.com",
				"password":"Pa55w0rd"
				}
			`,
			expect_error: true,
		},
		{
			test_name: "Test Other error",
			body: `{
				"name":"Other error",
				"email":"email@other.com",
				"password":"Pa55w0rd"
				}
			`,
			expect_error: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.test_name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(
				http.MethodPost,
				"/users",
				strings.NewReader(tt.body),
			)
			c.Request.Header.Set(
				"Content-Type",
				"application/json",
			)
			handler.Create(c)
			status := c.Writer.Status()

			if status == http.StatusCreated && tt.expect_error {
				t.Error("Error Expected:", status)
			}

			if status != http.StatusCreated && !tt.expect_error {
				t.Error("Unexpected error: ", status)
			}
		})
	}
}

func TestUserHandler_Get(t *testing.T) {
	userMock := &MockUserService{
		GetByEmailFn: func(ctx context.Context, email string) (*domain.User, error) {

			if email == EMAIL_NOT_FOUND {
				return nil, domain.ErrUserNotFound
			}

			if email == EMAIL_OTHER_ERROR {
				return nil, errors.New("Generic Error")
			}

			return &domain.User{
				ID:    ID_OK,
				Name:  NAME_OK,
				Email: EMAIL_OK,
			}, nil
		},
	}

	handler := NewUserHandler(userMock)

	tests := []struct {
		test_name    string
		email        string
		expect_error bool
	}{
		{
			test_name:    "Test OK",
			email:        EMAIL_OK,
			expect_error: false,
		},
		{
			test_name:    "Test Missing Email",
			expect_error: true,
		},
		{
			test_name:    "Test Email Not Found",
			email:        EMAIL_NOT_FOUND,
			expect_error: false,
		},
		{
			test_name:    "Test Generic Error",
			email:        EMAIL_OTHER_ERROR,
			expect_error: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.test_name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(
				http.MethodGet,
				fmt.Sprintf("/users?email=%s", tt.email),
				nil,
			)
			handler.Get(c)

			status := c.Writer.Status()
			if tt.expect_error && status == http.StatusOK {
				t.Errorf("Expecting error")
			}

		})
	}
}

func TestUserHandler_GetByID(t *testing.T) {

	mockUser := &MockUserService{
		GetByIdFn: func(ctx context.Context, id string) (*domain.User, error) {

			if id == ID_NOT_FOUND {
				return nil, domain.ErrUserNotFound
			}

			if id == "" {
				return nil, domain.ErrMissingAttribute
			}

			return &domain.User{
				ID:    ID_OK,
				Name:  NAME_OK,
				Email: EMAIL_OK,
			}, nil
		},
	}

	handler := NewUserHandler(mockUser)

	tests := []struct {
		test_name     string
		id            string
		expects_error bool
	}{
		{
			test_name:     "Test OK",
			id:            ID_OK,
			expects_error: false,
		},
		{
			test_name:     "Test missing ID",
			expects_error: true,
		},
		{
			test_name:     "Test Not Found",
			id:            ID_NOT_FOUND,
			expects_error: true,
		},
	}

	for _, tt := range tests {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/users/%s", tt.id),
			nil,
		)

		c.Params = gin.Params{
			{
				Key:   "id",
				Value: tt.id,
			},
		}

		handler.GetById(c)
		status := c.Writer.Status()
		if tt.expects_error && status == http.StatusOK {
			t.Errorf("Expecting error")
		}

	}
}

func TestUserHandler_Delete(t *testing.T) {
	mockUser := &MockUserService{
		DeleteFn: func(ctx context.Context, id string) error {

			if id == ID_NOT_FOUND {
				return domain.ErrUserNotFound
			}

			if id == "" {
				return domain.ErrMissingAttribute
			}
			if id == ID_GENERIC_ERROR {
				return errors.New("Generic Error")
			}

			return nil
		},
	}

	handler := NewUserHandler(mockUser)

	tests := []struct {
		test_name     string
		id            string
		expects_error bool
	}{
		{
			test_name:     "Test ok",
			id:            ID_OK,
			expects_error: false,
		},
		{
			test_name:     "ID not Found",
			id:            ID_NOT_FOUND,
			expects_error: true,
		},
		{
			test_name:     "ID not Found",
			id:            ID_GENERIC_ERROR,
			expects_error: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.test_name, func(t *testing.T) {

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodDelete,
				fmt.Sprintf("/users/%s", tt.id),
				nil)
			c.Params = gin.Params{
				{
					Key:   "id",
					Value: tt.id,
				},
			}

			handler.Delete(c)
			status := c.Writer.Status()
			if !tt.expects_error && status != http.StatusNoContent {
				t.Errorf("Error not expected")
			}

		})

	}

}

func TestUserHandler_Update(t *testing.T) {
	emailAlreadyExists := false
	genericError := false
	mockUser := &MockUserService{
		UpdateFn: func(ctx context.Context, id string, req servicedto.UpdateUserInput) (*domain.User, error) {
			if id == ID_NOT_FOUND {
				return nil, domain.ErrUserNotFound
			}
			if emailAlreadyExists {
				emailAlreadyExists = false
				return nil, domain.ErrEmailAlreadyExists
			}
			if genericError {
				genericError = false
				return nil, errors.New("Generic error")
			}
			return &domain.User{
				ID:        ID_OK,
				Name:      NAME_OK,
				Email:     EMAIL_OK,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now().Add(3 * time.Minute),
			}, nil
		},
	}

	handler := NewUserHandler(mockUser)

	tests := []struct {
		test_name     string
		id            string
		body          string
		expects_error bool
		email_exists  bool
		generic_error bool
	}{
		{
			test_name: "Test Ok",
			id:        ID_OK,
			body: `{
				"name":"Teste Ok",
				"email":"email@ok.com",
				"password":"Pa55w0rd"
				}
			`,
			expects_error: false,
			email_exists:  false,
			generic_error: false,
		},
		{
			test_name: "Test ID Not Found",
			id:        ID_NOT_FOUND,
			body: `{
				"name":"Teste Not Found",
				"email":"email@ok.com",
				"password":"Pa55w0rd"
				}
			`,
			expects_error: true,
			email_exists:  false,
			generic_error: false,
		},
		{
			test_name: "Test No Body",
			id:        ID_OK,
			body: `{

				}
			`,
			expects_error: true,
			email_exists:  false,
			generic_error: false,
		},
		{
			test_name: "Test Malformed Body",
			id:        ID_OK,
			body: `{
				malformed
				}
			`,
			expects_error: true,
			email_exists:  false,
			generic_error: false,
		},
		{
			test_name: "Test No ID",
			body: `{
				"name":"Teste Ok",
				"email":"email@ok.com",
				"password":"Pa55w0rd"
				}
			`,
			expects_error: true,
			email_exists:  false,
			generic_error: false,
		},
		{
			test_name: "Test update id",
			id:        ID_OK,
			body: `{
				"id":"update id not allowed",
				"name":"Teste Ok",
				"email":"email@ok.com",
				"password":"Pa55w0rd"
				}
			`,
			expects_error: true,
			email_exists:  false,
			generic_error: false,
		},
		{
			test_name: "Test Email Already Exists",
			id:        ID_OK,
			body: `{
				"name":"Teste Ok",
				"email":"email@already_exists.com",
				"password":"Pa55w0rd"
				}
			`,
			expects_error: true,
			email_exists:  true,
			generic_error: false,
		},
		{
			test_name: "Generic Error",
			id:        ID_OK,
			body: `{
				"name":"Teste Ok",
				"email":"email@already_exists.com",
				"password":"Pa55w0rd"
				}
			`,
			expects_error: true,
			email_exists:  false,
			generic_error: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.test_name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			if tt.email_exists {
				emailAlreadyExists = true
			}
			if tt.generic_error {
				genericError = true
			}
			c.Request = httptest.NewRequest(
				http.MethodPatch,
				fmt.Sprintf("/users/%s", tt.id),
				strings.NewReader(tt.body),
			)
			c.Params = gin.Params{
				{
					Key:   "id",
					Value: tt.id,
				},
			}

			handler.Update(c)
			status := c.Writer.Status()

			if status != 200 && !tt.expects_error {
				t.Errorf("Unexpected error")
			}

		})
	}
}
