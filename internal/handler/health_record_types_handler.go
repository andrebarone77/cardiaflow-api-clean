package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/andrebarone77/cardiaflow-api/internal/domain"
	handlerdto "github.com/andrebarone77/cardiaflow-api/internal/handler/dto"
	servicedto "github.com/andrebarone77/cardiaflow-api/internal/service/dto"
	"github.com/gin-gonic/gin"
)

type HealthRecordTypeService interface {
	Create(ctx context.Context, input servicedto.HealthRecordTypeInput) (string, error)
	GetByID(ctx context.Context, id string) (*domain.HealthRecordType, error)
	GetByCode(ctx context.Context, code string) (*domain.HealthRecordType, error)
	GetAll(ctx context.Context) ([]*domain.HealthRecordType, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, update servicedto.HealthRecordTypeUpdateInput) (*domain.HealthRecordType, error)
}

type HealthRecordTypeHandler struct {
	healthRecordTypeService HealthRecordTypeService
}

func NewHealthRecordTypeHandler(service HealthRecordTypeService) *HealthRecordTypeHandler {
	return &HealthRecordTypeHandler{healthRecordTypeService: service}
}

func (h *HealthRecordTypeHandler) Create(c *gin.Context) {
	var req handlerdto.CreateHealthRecordTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.healthRecordTypeService.Create(c.Request.Context(), toHealthTypeServiceCreateInput(req))

	if err != nil {
		if errors.Is(err, domain.ErrHealthRecordTypeAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"message": domain.ErrHealthRecordTypeAlreadyExists.Error()})
			return
		}

		if errors.Is(err, domain.ErrCodeRequired) {
			c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrCodeRequired.Error()})
			return
		}

		if errors.Is(err, domain.ErrCodeInvalid) {
			c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrCodeInvalid.Error()})
			return
		}

		if errors.Is(err, domain.ErrCodeTooLong) {
			c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrCodeTooLong.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid request body"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *HealthRecordTypeHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	healthRecordType, err := h.healthRecordTypeService.GetByID(c, id)
	if err != nil {
		if errors.Is(err, domain.ErrHealthRecordTypeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message:": domain.ErrHealthRecordTypeNotFound.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":   healthRecordType.ID,
		"name": healthRecordType.Name,
		"code": healthRecordType.Code,
	})
}

func (h *HealthRecordTypeHandler) GetByCode(c *gin.Context) {
	code := c.Param("code")

	healthRecordType, err := h.healthRecordTypeService.GetByCode(c, code)
	if err != nil {
		if errors.Is(err, domain.ErrHealthRecordTypeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message:": domain.ErrHealthRecordTypeNotFound.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":   healthRecordType.ID,
		"name": healthRecordType.Name,
		"code": healthRecordType.Code,
	})
}

func (h *HealthRecordTypeHandler) GetAll(c *gin.Context) {
	healthRecordTypes, err := h.healthRecordTypeService.GetAll(c)

	if err != nil {
		if errors.Is(err, domain.ErrHealthRecordTypeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message:": domain.ErrHealthRecordTypeNotFound.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, healthRecordTypes)
}

func (h *HealthRecordTypeHandler) Delete(c *gin.Context) {
	id := c.Query("id")

	err := h.healthRecordTypeService.Delete(c, id)

	if err != nil {
		if errors.Is(err, domain.ErrHealthRecordTypeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": domain.ErrHealthRecordTypeNotFound.Error()})
			return
		}

		if errors.Is(err, domain.ErrHealthRecordTypeImmutable) {
			c.JSON(http.StatusForbidden, gin.H{"error": domain.ErrHealthRecordTypeImmutable.Error()})
			return
		}
	}

	c.Status(http.StatusNoContent)

}

func (h *HealthRecordTypeHandler) Update(c *gin.Context) {
	var req handlerdto.UpdateHealthRecordTypeRequest

	id := c.Param("id")

	if len(id) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrEmptyId.Error()})
		return
	}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name == nil && req.Code == nil && req.Unit == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrNoInformation.Error()})
		return
	}

	healthRecordType, err := h.healthRecordTypeService.Update(c.Request.Context(), id, toHealthTypeServiceCreateUpdate(req))

	if err != nil {
		if errors.Is(err, domain.ErrHealthRecordTypeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": domain.ErrHealthRecordTypeNotFound.Error()})
			return
		}

		if errors.Is(err, domain.ErrHealthRecordTypeAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": domain.ErrCodeAlreadyExists.Error()})
			return
		}

		if errors.Is(err, domain.ErrHealthRecordTypeImmutable) {
			c.JSON(http.StatusForbidden, gin.H{"error": domain.ErrHealthRecordTypeImmutable.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name": healthRecordType.Name,
		"code": healthRecordType.Code,
		"unit": healthRecordType.Unit,
		"id":   healthRecordType.ID,
	})

}

func toHealthTypeServiceCreateInput(req handlerdto.CreateHealthRecordTypeRequest) servicedto.HealthRecordTypeInput {
	return servicedto.HealthRecordTypeInput{
		Name: req.Name,
		Code: req.Code,
		Unit: req.Unit,
	}
}
func toHealthTypeServiceCreateUpdate(req handlerdto.UpdateHealthRecordTypeRequest) servicedto.HealthRecordTypeUpdateInput {
	return servicedto.HealthRecordTypeUpdateInput{
		Name: req.Name,
		Code: req.Code,
		Unit: req.Unit,
	}
}
