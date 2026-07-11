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

// Login godoc
//
// @Summary      Authenticate user
// @Description  Authenticates a user and returns a JWT token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginRequest true "Login credentials"
// @Success      200 {object} dto.LoginResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      400 {object} dto.ErrorResponse
// @Router       /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req handlerdto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse := handlerdto.ErrorResponse{
			Error: err.Error(),
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	token_resp, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)

	if err != nil {
		errorResponse := handlerdto.ErrorResponse{
			Error: domain.ErrNotAuthorized.Error(),
		}
		c.JSON(http.StatusUnauthorized, errorResponse)
		return
	}

	loginResponse := handlerdto.LoginResponse{
		Token: token_resp,
	}

	c.JSON(http.StatusOK, loginResponse)
}
