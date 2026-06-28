package handler

import (
	"context"
	"net/http"

	"github.com/andrebarone77/cardiaflow-api/internal/domain"
	handlerdto "github.com/andrebarone77/cardiaflow-api/internal/handler/dto"
	"github.com/gin-gonic/gin"
)

type AuthService interface {
	Login(ctx context.Context, email string, password string) (string, error)
}

type AuthHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req handlerdto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": domain.ErrNotAuthorized.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
