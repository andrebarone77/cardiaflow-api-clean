package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/andrebarone77/cardiaflow-api/internal/domain"
	handlerdto "github.com/andrebarone77/cardiaflow-api/internal/handler/dto"
	"github.com/andrebarone77/cardiaflow-api/internal/service/dto"
	servicedto "github.com/andrebarone77/cardiaflow-api/internal/service/dto"
)

type UserService interface {
	Create(ctx context.Context, req servicedto.CreateUserInput) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetById(ctx context.Context, id string) (*domain.User, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, req servicedto.UpdateUserInput) (*domain.User, error)
}

type UserHandler struct {
	userService UserService
}

func NewUserHandler(userService UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Create godoc
//
// @Summary      Creates User
// @Description  Returns created user
// @Tags         Users
// @Security     BearerAuth
// @Produce      json
// @Param        request body dto.CreateUserRequest true "User data"
// @Success      201 {object} dto.UserResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      409 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /api/users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req handlerdto.CreateUserRequest

	//1. Bind + validação básica
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, handlerdto.ErrorResponse{
			Error: err.Error()})
		return
	}

	//2. Chama o service
	user, err := h.userService.Create(c.Request.Context(), toServiceCreateInput(req))
	if err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyExists) {
			c.JSON(http.StatusConflict,
				handlerdto.ErrorResponse{Error: domain.ErrEmailAlreadyExists.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError,
			handlerdto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, handlerdto.UserResponse{ID: user.ID,
		Name:  user.Name,
		Email: user.Email})

}

// Get godoc
//
// @Summary      Get user by E-mail
// @Description  Returns a user from its Email
// @Tags         Users
// @Security     BearerAuth
// @Produce      json
// @Param        email query string true "User e-mail"
// @Success      200 {object} dto.UserResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /api/users [get]
func (h *UserHandler) Get(c *gin.Context) {
	email := strings.ToLower(strings.TrimSpace(c.Query("email")))

	if email == "" {
		c.JSON(http.StatusBadRequest,
			handlerdto.ErrorResponse{Error: "Missing email"})
		return
	}
	user, err := h.userService.GetByEmail(c.Request.Context(), email)

	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, handlerdto.ErrorResponse{
				Error: domain.ErrUserNotFound.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError,
			handlerdto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, handlerdto.UserResponse{ID: user.ID,
		Name:  user.Name,
		Email: user.Email})

}

// GetById godoc
//
// @Summary      Get user by ID
// @Description  Returns a user from its UUID
// @Tags         Users
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "User UUID"
// @Success      200 {object} dto.UserResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /api/users/{id} [get]
func (h *UserHandler) GetById(c *gin.Context) {
	id := c.Param("id")

	user, err := h.userService.GetById(c.Request.Context(), id)

	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, handlerdto.ErrorResponse{
				Error: domain.ErrUserNotFound.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, handlerdto.ErrorResponse{
			Error: err.Error(),
		})

		return

	}

	response := handlerdto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	c.JSON(http.StatusOK, response)

}

// Delete godoc
//
// @Summary      Delete user by ID
// @Description  Delete a User by its UUID
// @Tags         Users
// @Security     BearerAuth
// @Param        id path string true "User UUID"
// @Success      204
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /api/users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.userService.Delete(c.Request.Context(), id)

	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, handlerdto.ErrorResponse{
				Error: domain.ErrUserNotFound.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, handlerdto.ErrorResponse{
			Error: err.Error(),
		})

		return

	}

	c.Status(http.StatusNoContent)
}

func (h *UserHandler) Update(c *gin.Context) {
	var req handlerdto.UpdateUserRequest

	id := c.Param("id")

	if len(id) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrEmptyId.Error()})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.ID != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "update of ID is not allowed"})
		return
	}

	if req.Email == nil && req.Name == nil && req.Password == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrNoInformation.Error()})
		return
	}

	user, err := h.userService.Update(c.Request.Context(), id, toServiceUpdateInput(req))

	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": domain.ErrUserNotFound.Error()})
			return
		}

		if errors.Is(err, domain.ErrEmailAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": domain.ErrEmailAlreadyExists.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": user.ID,
		"name":  user.Name,
		"email": user.Email})

}

func toServiceCreateInput(req handlerdto.CreateUserRequest) dto.CreateUserInput {
	return servicedto.CreateUserInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
}

func toServiceUpdateInput(req handlerdto.UpdateUserRequest) servicedto.UpdateUserInput {
	return servicedto.UpdateUserInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
}
