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

// Create godoc
//
// @Summary      Create Health Record Type
// @Description  Creates a new Health Record Type
// @Tags         Health Record Type
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateHealthRecordTypeRequest true "Health Record Type data"
// @Success      201 {object} dto.CreateHealthRecordTypeResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      409 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /api/healthrecordtypes [post]
func (h *HealthRecordTypeHandler) Create(c *gin.Context) {
	var req handlerdto.CreateHealthRecordTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, handlerdto.ErrorResponse{Error: err.Error()})
		return
	}

	id, err := h.healthRecordTypeService.Create(c.Request.Context(), toHealthTypeServiceCreateInput(req))

	if err != nil {
		if errors.Is(err, domain.ErrHealthRecordTypeAlreadyExists) {
			c.JSON(http.StatusConflict, handlerdto.ErrorResponse{Error: domain.ErrHealthRecordTypeAlreadyExists.Error()})
			return
		}

		if errors.Is(err, domain.ErrCodeRequired) {
			c.JSON(http.StatusBadRequest, handlerdto.ErrorResponse{Error: domain.ErrCodeRequired.Error()})
			return
		}

		if errors.Is(err, domain.ErrCodeInvalid) {
			c.JSON(http.StatusBadRequest, handlerdto.ErrorResponse{Error: domain.ErrCodeInvalid.Error()})
			return
		}

		if errors.Is(err, domain.ErrCodeTooLong) {
			c.JSON(http.StatusBadRequest, handlerdto.ErrorResponse{Error: domain.ErrCodeTooLong.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, handlerdto.ErrorResponse{Error: "invalid request body"})
		return
	}
	response := handlerdto.CreateHealthRecordTypeResponse{
		Id: id,
	}
	c.JSON(http.StatusCreated, response)
}

// GetById godoc
//
// @Summary      Get health record type by ID
// @Description  Returns a health record type from its UUID
// @Tags         Health Record Type
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "Health Record Type UUID"
// @Success      200 {object} dto.HealthRecordTypeResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /api/healthrecordtypes/{id} [get]
func (h *HealthRecordTypeHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	healthRecordType, err := h.healthRecordTypeService.GetByID(c, id)
	if err != nil {
		if errors.Is(err, domain.ErrHealthRecordTypeNotFound) {
			c.JSON(http.StatusNotFound, handlerdto.ErrorResponse{Error: domain.ErrHealthRecordTypeNotFound.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, handlerdto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, handlerdto.HealthRecordTypeResponse{
		Id:   healthRecordType.ID,
		Name: healthRecordType.Name,
		Code: healthRecordType.Code,
	})
}

// GetByCode godoc
//
// @Summary      Get health record type by Code
// @Description  Returns a health record type from its code
// @Tags         Health Record Type
// @Security     BearerAuth
// @Produce      json
// @Param        code path string true "Health Record Type UUID"
// @Success      200 {object} dto.HealthRecordTypeResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /api/healthrecordtypes/code/{code} [get]
func (h *HealthRecordTypeHandler) GetByCode(c *gin.Context) {
	code := c.Param("code")

	healthRecordType, err := h.healthRecordTypeService.GetByCode(c, code)
	if err != nil {
		if errors.Is(err, domain.ErrHealthRecordTypeNotFound) {
			c.JSON(http.StatusNotFound, handlerdto.ErrorResponse{Error: domain.ErrHealthRecordTypeNotFound.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, handlerdto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, handlerdto.HealthRecordTypeResponse{
		Id:   healthRecordType.ID,
		Name: healthRecordType.Name,
		Code: healthRecordType.Code,
	})
}

// GetAll godoc
//
// @Summary      Get all health record types
// @Description  Returns a health record type from its code
// @Tags         Health Record Type
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} []dto.HealthRecordTypeResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /api/healthrecordtypes [Get]
func (h *HealthRecordTypeHandler) GetAll(c *gin.Context) {
	healthRecordTypes, err := h.healthRecordTypeService.GetAll(c)

	if err != nil {
		if errors.Is(err, domain.ErrHealthRecordTypeNotFound) {
			c.JSON(http.StatusNotFound, handlerdto.ErrorResponse{Error: domain.ErrHealthRecordTypeNotFound.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, handlerdto.ErrorResponse{Error: err.Error()})
		return
	}

	healthRecordTypesResponse := toHealthRecordTypeResponseList(healthRecordTypes)
	c.JSON(http.StatusOK, healthRecordTypesResponse)
}

// Delete godoc
//
// @Summary      Delete a Health Record Type by ID
// @Description  Delete a Health Record Type by ID
// @Tags         Health Record Type
// @Security     BearerAuth
// @Produce      json
// @Param        code query string true "Health Record Type UUID"
// @Success      209
// @Failure      404 {object} dto.ErrorResponse
// @Failure      403 {object} dto.ErrorResponse
// @Router       /api/healthrecordtypes{id} [delete]
func (h *HealthRecordTypeHandler) Delete(c *gin.Context) {
	id := c.Query("id")

	err := h.healthRecordTypeService.Delete(c, id)

	if err != nil {
		if errors.Is(err, domain.ErrHealthRecordTypeNotFound) {
			c.JSON(http.StatusNotFound, handlerdto.ErrorResponse{Error: domain.ErrHealthRecordTypeNotFound.Error()})
			return
		}

		if errors.Is(err, domain.ErrHealthRecordTypeImmutable) {
			c.JSON(http.StatusForbidden, handlerdto.ErrorResponse{Error: domain.ErrHealthRecordTypeImmutable.Error()})
			return
		}
	}

	c.Status(http.StatusNoContent)

}

// Update godoc
//
// @Summary      Update a Health Record Type by ID
// @Description  Update a Health Record Type by ID
// @Tags         Health Record Type
// @Security     BearerAuth
// @Accept		 json
// @Produce      json
// @Param        id path string true "Health Record Type UUID"
// @Param		 request body dto.UpdateHealthRecordTypeRequest true "Health Record Type data"
// @Success      200 {object} dto.HealthRecordTypeUpdateResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      403 {object} dto.ErrorResponse
// @Failure      409 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /api/healthrecordtypes/{id} [patch]
func (h *HealthRecordTypeHandler) Update(c *gin.Context) {
	var req handlerdto.UpdateHealthRecordTypeRequest

	id := c.Param("id")

	if len(id) == 0 {
		c.JSON(http.StatusBadRequest, handlerdto.ErrorResponse{Error: domain.ErrEmptyId.Error()})
		return
	}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, handlerdto.ErrorResponse{Error: err.Error()})
		return
	}

	if req.Name == nil && req.Code == nil && req.Unit == nil {
		c.JSON(http.StatusBadRequest, handlerdto.ErrorResponse{Error: domain.ErrNoInformation.Error()})
		return
	}

	healthRecordType, err := h.healthRecordTypeService.Update(c.Request.Context(), id, toHealthTypeServiceCreateUpdate(req))

	if err != nil {
		if errors.Is(err, domain.ErrHealthRecordTypeNotFound) {
			c.JSON(http.StatusNotFound, handlerdto.ErrorResponse{Error: domain.ErrHealthRecordTypeNotFound.Error()})
			return
		}

		if errors.Is(err, domain.ErrHealthRecordTypeAlreadyExists) {
			c.JSON(http.StatusConflict, handlerdto.ErrorResponse{Error: domain.ErrCodeAlreadyExists.Error()})
			return
		}

		if errors.Is(err, domain.ErrHealthRecordTypeImmutable) {
			c.JSON(http.StatusForbidden, handlerdto.ErrorResponse{Error: domain.ErrHealthRecordTypeImmutable.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, handlerdto.ErrorResponse{Error: "internal error"})
		return
	}

	c.JSON(http.StatusOK, handlerdto.HealthRecordTypeUpdateResponse{
		Name: healthRecordType.Name,
		Code: healthRecordType.Code,
		Unit: *healthRecordType.Unit,
		Id:   healthRecordType.ID,
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

func toHealthRecordTypeResponseList(healthRecordTypes []*domain.HealthRecordType) []handlerdto.HealthRecordTypeResponse {
	var healthRecordTypesResponse []handlerdto.HealthRecordTypeResponse

	for _, healthRecordTypeResponse := range healthRecordTypes {
		hrt := handlerdto.HealthRecordTypeResponse{
			Id:   healthRecordTypeResponse.ID,
			Name: healthRecordTypeResponse.Name,
			Code: healthRecordTypeResponse.Code,
		}
		healthRecordTypesResponse = append(healthRecordTypesResponse, hrt)
	}

	return healthRecordTypesResponse

}
