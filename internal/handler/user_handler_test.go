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

type MockUserService struct {
	CreateFn     func(ctx context.Context, req servicedto.CreateUserInput) (*domain.User, error)
	GetByEmailFn func(ctx context.Context, email string, userID string, requesterRole domain.Role) (*domain.User, error)
	GetByIdFn    func(ctx context.Context, id string, userID string, requesterRole domain.Role) (*domain.User, error)
	DeleteFn     func(ctx context.Context, id string, userID string, requesterRole domain.Role) error
	UpdateFn     func(ctx context.Context, id string, userID string, requesterRole domain.Role, req servicedto.UpdateUserInput) (*domain.User, error)
}

func (m *MockUserService) Create(ctx context.Context, req servicedto.CreateUserInput) (*domain.User, error) {
	return m.CreateFn(ctx, req)
}

func (m *MockUserService) GetByEmail(ctx context.Context, email string, userID string, requesterRole domain.Role) (*domain.User, error) {
	return m.GetByEmailFn(ctx, email, userID, requesterRole)
}

func (m *MockUserService) GetById(ctx context.Context, id string, userID string, requesterRole domain.Role) (*domain.User, error) {
	return m.GetByIdFn(ctx, id, userID, requesterRole)
}

func (m *MockUserService) Delete(ctx context.Context, id string, userID string, requesterRole domain.Role) error {
	return m.DeleteFn(ctx, id, userID, requesterRole)
}

func (m *MockUserService) Update(ctx context.Context, id string, userID string, requesterRole domain.Role, req servicedto.UpdateUserInput) (*domain.User, error) {
	return m.UpdateFn(ctx, id, userID, requesterRole, req)
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
		GetByEmailFn: func(ctx context.Context, email string, userID string, requesterRole domain.Role) (*domain.User, error) {

			if email == EMAIL_NOT_FOUND {
				return nil, domain.ErrUserNotFound
			}

			if email == EMAIL_NOT_AUTHORIZED && requesterRole == domain.RoleUser {
				return nil, domain.ErrForbidden
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
		test_name          string
		email              string
		expect_error       bool
		requester_id       string
		requester_num      int8
		requester_role     domain.Role
		requester_role_num int8
	}{
		{
			test_name:          "Test OK",
			email:              EMAIL_OK,
			expect_error:       false,
			requester_id:       USER_ID_REGULAR,
			requester_num:      0,
			requester_role:     domain.RoleUser,
			requester_role_num: 0,
		},
		{
			test_name:          "Test Missing Email",
			expect_error:       true,
			requester_id:       USER_ID_REGULAR,
			requester_num:      0,
			requester_role:     domain.RoleUser,
			requester_role_num: 0,
		},
		{
			test_name:          "Test Email Not Found",
			email:              EMAIL_NOT_FOUND,
			expect_error:       true,
			requester_id:       USER_ID_REGULAR,
			requester_num:      0,
			requester_role:     domain.RoleUser,
			requester_role_num: 0,
		},
		{
			test_name:          "Test Email Forbidden",
			email:              EMAIL_NOT_AUTHORIZED,
			expect_error:       true,
			requester_id:       USER_ID_REGULAR,
			requester_num:      0,
			requester_role:     domain.RoleUser,
			requester_role_num: 0,
		},
		{
			test_name:          "Test Generic Error",
			email:              EMAIL_OTHER_ERROR,
			expect_error:       true,
			requester_id:       USER_ID_REGULAR,
			requester_num:      0,
			requester_role:     domain.RoleUser,
			requester_role_num: 0,
		},
		{
			test_name:          "Test RequesterID Conversion Error",
			email:              EMAIL_OK,
			expect_error:       true,
			requester_id:       USER_ID_REGULAR,
			requester_num:      123,
			requester_role:     domain.RoleUser,
			requester_role_num: 0,
		},
		{
			test_name:          "User Forbidden",
			email:              EMAIL_NOT_AUTHORIZED,
			expect_error:       true,
			requester_id:       USER_ID_REGULAR,
			requester_num:      123,
			requester_role:     domain.RoleUser,
			requester_role_num: 0,
		},
		{
			test_name:          "Missing Requester ID",
			email:              EMAIL_OK,
			expect_error:       true,
			requester_num:      0,
			requester_role:     domain.RoleUser,
			requester_role_num: 0,
		},
		{
			test_name:          "Missing Requester Role",
			email:              EMAIL_OK,
			expect_error:       true,
			requester_id:       USER_ID_REGULAR,
			requester_num:      0,
			requester_role_num: 0,
		},
		{
			test_name:          "Test Role Conversion Error",
			email:              EMAIL_OK,
			expect_error:       true,
			requester_id:       USER_ID_REGULAR,
			requester_num:      0,
			requester_role_num: 123,
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
			if tt.requester_num == 0 && tt.requester_id != "" {
				c.Set("userID", tt.requester_id)
			}
			if tt.requester_num == 123 {
				c.Set("userID", tt.requester_num)
			}

			if tt.requester_role != "" && tt.requester_role_num != 123 {
				c.Set("role", tt.requester_role)
			}

			if tt.requester_role_num == 123 {
				c.Set("role", tt.requester_role_num)
			}

			handler.Get(c)

			status := c.Writer.Status()
			if tt.expect_error && status == http.StatusOK {
				t.Errorf("Expecting error")
			}

			if !tt.expect_error && status != http.StatusOK {
				t.Errorf(("Unexptected Error"))
			}

		})
	}
}

func TestUserHandler_GetByID(t *testing.T) {

	mockUser := &MockUserService{
		GetByIdFn: func(ctx context.Context, id string, userID string, requesterRole domain.Role) (*domain.User, error) {

			if id == ID_NOT_FOUND {
				return nil, domain.ErrUserNotFound
			}

			if id == "" {
				return nil, domain.ErrMissingAttribute
			}

			if id == ID_GENERIC_ERROR {
				return nil, errors.New("Generic Error")
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
		test_name          string
		id                 string
		requester_id       string
		requester_num      int8
		requester_role     domain.Role
		requester_role_num int8
		expects_error      bool
	}{
		{
			test_name:          "Test OK",
			id:                 ID_OK,
			requester_id:       USER_ID_REGULAR,
			requester_num:      0,
			requester_role:     domain.RoleUser,
			requester_role_num: 0,
			expects_error:      false,
		},
		{
			test_name:          "Test missing ID",
			requester_num:      0,
			requester_role:     domain.RoleUser,
			requester_role_num: 0,
			expects_error:      true,
		},
		{
			test_name:          "Test Not Found",
			id:                 ID_NOT_FOUND,
			requester_id:       USER_ID_REGULAR,
			requester_num:      0,
			requester_role:     domain.RoleUser,
			requester_role_num: 0,
			expects_error:      true,
		},
		{
			test_name:          "Test Generic Error",
			id:                 ID_GENERIC_ERROR,
			requester_id:       USER_ID_REGULAR,
			requester_num:      0,
			requester_role:     domain.RoleUser,
			requester_role_num: 0,
			expects_error:      true,
		},
		{
			test_name:          "Test Requester Fail",
			id:                 ID_OK,
			requester_id:       USER_ID_REGULAR,
			requester_num:      123,
			requester_role:     domain.RoleUser,
			requester_role_num: 0,
			expects_error:      true,
		},
		{
			test_name:          "Test Role Fail",
			id:                 ID_OK,
			requester_id:       USER_ID_REGULAR,
			requester_num:      0,
			requester_role:     domain.RoleUser,
			requester_role_num: 123,
			expects_error:      true,
		},
		{
			test_name:          "Test Requester Empty",
			id:                 ID_OK,
			requester_role:     domain.RoleUser,
			requester_role_num: 0,
			expects_error:      true,
		},
		{
			test_name:     "Test Role Fail",
			id:            ID_OK,
			requester_id:  USER_ID_REGULAR,
			requester_num: 0,
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

		if tt.requester_id != "" && tt.requester_num == 0 {
			c.Set("userID", tt.requester_id)
		}

		if tt.requester_num == 123 {
			c.Set("userID", tt.requester_num)
		}

		if tt.requester_role != "" && tt.requester_role_num == 0 {
			c.Set("role", tt.requester_role)
		}

		if tt.requester_role_num == 123 {
			c.Set("role", tt.requester_role_num)
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
		DeleteFn: func(ctx context.Context, id string, userID string, requesterRole domain.Role) error {

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
		test_name          string
		id                 string
		requester_role     domain.Role
		requester_role_num int8
		requester_id       string
		requester_num      int8
		expects_error      bool
	}{
		{
			test_name:          "Test ok - User Admin",
			id:                 ID_OK,
			requester_role:     domain.RoleAdmin,
			requester_role_num: 0,
			requester_id:       USER_ID_ADMIN,
			requester_num:      0,
			expects_error:      false,
		},
		{
			test_name:          "Test ok - User Manager",
			id:                 ID_OK,
			requester_role:     domain.RoleAdmin,
			requester_role_num: 0,
			requester_id:       USER_ID_MANAGER,
			requester_num:      0,
			expects_error:      false,
		},
		{
			test_name:          "ID not Found",
			id:                 ID_NOT_FOUND,
			requester_role:     domain.RoleAdmin,
			requester_role_num: 0,
			requester_id:       USER_ID_ADMIN,
			requester_num:      0,
			expects_error:      true,
		},
		{
			test_name:          "ID not Found",
			id:                 ID_NOT_FOUND,
			requester_role:     domain.RoleUser,
			requester_role_num: 0,
			requester_id:       USER_ID_ADMIN,
			requester_num:      0,
			expects_error:      true,
		},
		{
			test_name:          "Generic error",
			id:                 ID_GENERIC_ERROR,
			requester_role:     domain.RoleAdmin,
			requester_role_num: 0,
			requester_id:       USER_ID_ADMIN,
			requester_num:      0,
			expects_error:      true,
		},
		{
			test_name:          "Missing UserID",
			id:                 ID_OK,
			requester_role:     domain.RoleAdmin,
			requester_role_num: 0,
			requester_num:      0,
			expects_error:      true,
		},
		{
			test_name:          "Missing Role",
			id:                 ID_OK,
			requester_id:       USER_ID_ADMIN,
			requester_num:      0,
			requester_role_num: 0,
			expects_error:      true,
		},
		{
			test_name:          "Generic error",
			id:                 ID_OK,
			requester_role:     domain.RoleAdmin,
			requester_role_num: 0,
			requester_id:       USER_ID_ADMIN,
			requester_num:      0,
			expects_error:      true,
		},
		{
			test_name:          "Generic error",
			id:                 ID_GENERIC_ERROR,
			requester_role:     domain.RoleAdmin,
			requester_role_num: 0,
			requester_id:       USER_ID_ADMIN,
			requester_num:      0,
			expects_error:      true,
		},
		{
			test_name:          "Requester Missing",
			id:                 ID_GENERIC_ERROR,
			requester_role:     domain.RoleAdmin,
			requester_role_num: 0,
			requester_num:      0,
			expects_error:      true,
		},
		{
			test_name:          "Role Missing",
			id:                 ID_GENERIC_ERROR,
			requester_role_num: 0,
			requester_id:       USER_ID_ADMIN,
			requester_num:      0,
			expects_error:      true,
		},
		{
			test_name:          "Invalid Requester",
			id:                 ID_GENERIC_ERROR,
			requester_role:     domain.RoleAdmin,
			requester_role_num: 0,
			requester_num:      123,
			expects_error:      true,
		},
		{
			test_name:          "Invalid Role",
			id:                 ID_GENERIC_ERROR,
			requester_role_num: 123,
			requester_id:       USER_ID_ADMIN,
			requester_num:      0,
			expects_error:      true,
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

			if tt.requester_id != "" && tt.requester_num != 123 {
				c.Set("userID", tt.requester_id)
			}

			if tt.requester_num == 123 {
				c.Set("userID", tt.requester_num)
			}

			if tt.requester_role != "" && tt.requester_role_num != 123 {
				c.Set("role", tt.requester_role)
			}

			if tt.requester_role_num == 123 {
				c.Set("role", tt.requester_role_num)
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
		UpdateFn: func(ctx context.Context, id string, userID string, requesterRole domain.Role, req servicedto.UpdateUserInput) (*domain.User, error) {
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
		requester_id  string
		requester_num int8
		role          domain.Role
		role_num      int8
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
			requester_id:  USER_ID_REGULAR,
			requester_num: 0,
			role:          domain.RoleUser,
			role_num:      0,
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
			requester_id:  USER_ID_REGULAR,
			requester_num: 0,
			role:          domain.RoleUser,
			role_num:      0,
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
			requester_id:  USER_ID_REGULAR,
			requester_num: 0,
			role:          domain.RoleUser,
			role_num:      0,
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
			requester_id:  USER_ID_REGULAR,
			requester_num: 0,
			role:          domain.RoleUser,
			role_num:      0,
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
			requester_id:  USER_ID_REGULAR,
			requester_num: 0,
			role:          domain.RoleUser,
			role_num:      0,
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
			requester_id:  USER_ID_REGULAR,
			requester_num: 0,
			role:          domain.RoleUser,
			role_num:      0,
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
			requester_id:  USER_ID_REGULAR,
			requester_num: 0,
			role:          domain.RoleUser,
			role_num:      0,
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
			requester_id:  USER_ID_REGULAR,
			requester_num: 0,
			role:          domain.RoleUser,
			role_num:      0,
		},
		{
			test_name: "Test Missing Requester ID",
			id:        ID_OK,
			body: `{
				"name":"Teste Ok",
				"email":"email@ok.com",
				"password":"Pa55w0rd"
				}
			`,
			expects_error: true,
			email_exists:  false,
			generic_error: false,
			role:          domain.RoleUser,
			role_num:      0,
		},
		{
			test_name: "Test Missing Requester Role",
			id:        ID_OK,
			body: `{
				"name":"Teste Ok",
				"email":"email@ok.com",
				"password":"Pa55w0rd"
				}
			`,
			expects_error: true,
			email_exists:  false,
			generic_error: false,
			requester_id:  USER_ID_REGULAR,
			requester_num: 0,
		},
		{
			test_name: "Test Malformed Requester ID",
			id:        ID_OK,
			body: `{
				"name":"Teste Ok",
				"email":"email@ok.com",
				"password":"Pa55w0rd"
				}
			`,
			expects_error: true,
			email_exists:  false,
			generic_error: false,
			requester_num: 123,
			role:          domain.RoleUser,
			role_num:      0,
		},
		{
			test_name: "Test Malformed Role",
			id:        ID_OK,
			body: `{
				"name":"Teste Ok",
				"email":"email@ok.com",
				"password":"Pa55w0rd"
				}
			`,
			expects_error: true,
			email_exists:  false,
			generic_error: false,
			requester_id:  USER_ID_REGULAR,
			requester_num: 0,
			role_num:      123,
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

			if tt.requester_id != "" && tt.requester_num == 0 {
				c.Set("userID", tt.requester_id)
			}

			if tt.requester_num == 123 {
				c.Set("userID", tt.requester_num)
			}

			if tt.role_num == 0 && tt.role != "" {
				c.Set("role", tt.role)
			}

			if tt.role_num == 123 {
				c.Set("role", tt.role_num)
			}

			handler.Update(c)
			status := c.Writer.Status()

			if status != http.StatusOK && !tt.expects_error {
				t.Errorf("Unexpected error")
			}

		})
	}
}
