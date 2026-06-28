package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAuthMiddleware_Uanuthorized(t *testing.T) {

	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(AuthMiddleware())

	r.GET("/protected", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(
		http.MethodGet,
		"/protected",
		nil,
	)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected %d, got %d",
			http.StatusUnauthorized,
			w.Code)
	}
}

func TestAuthMiddleware_InvalidAuthorization(t *testing.T) {

	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(AuthMiddleware())

	r.GET("/protected", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(
		http.MethodGet,
		"/protected",
		nil,
	)

	req.Header.Set("Authorization", "Basic abc123")

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected %d, got %d",
			http.StatusUnauthorized,
			w.Code)
	}
}

func TestAuthMiddleware_Authorized(t *testing.T) {

	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(AuthMiddleware())

	r.GET("/protected", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(
		http.MethodGet,
		"/protected",
		nil,
	)

	req.Header.Set("Authorization", "Bearer abc123")

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected %d, got %d",
			http.StatusUnauthorized,
			w.Code)
	}
}
