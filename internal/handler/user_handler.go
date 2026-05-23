package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/andrebarone77/cardiaflow-api/internal/domain"
	handlerdto "github.com/andrebarone77/cardiaflow-api/internal/handler/dto"
	"github.com/andrebarone77/cardiaflow-api/internal/service"
	"github.com/andrebarone77/cardiaflow-api/internal/service/dto"
	servicedto "github.com/andrebarone77/cardiaflow-api/internal/service/dto"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Create(c *gin.Context) {
	var req handlerdto.CreateUserRequest

	//1. Bind + validação básica
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//2. Chama o service
	user, err := h.userService.Create(c.Request.Context(), toServiceCreateInput(req))
	if err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": domain.ErrEmailAlreadyExists.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": user.ID,
		"name":  user.Name,
		"email": user.Email})

}

func (h *UserHandler) Get(c *gin.Context) {
	email := strings.ToLower(strings.TrimSpace(c.Query("email")))

	user, err := h.userService.Get(c.Request.Context(), email)

	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": domain.ErrUserNotFound.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": user.ID,
		"name":  user.Name,
		"email": user.Email})

}

func (h *UserHandler) GetById(c *gin.Context) {
	id := c.Param("id")

	user, err := h.userService.GetById(c.Request.Context(), id)

	if err != nil {
		if errors.Is(err, domain.ErrHealthRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": domain.ErrHealthRecordNotFound.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": user.ID,
		"name":  user.Name,
		"email": user.Email})

}

func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.userService.Delete(c.Request.Context(), id)

	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": domain.ErrUserNotFound.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
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
