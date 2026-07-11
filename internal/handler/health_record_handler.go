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

type HealthRecordService interface {
	Create(ctx context.Context, healthRecordInput servicedto.HealthRecordCreateInput) (string, error)
	GetByID(ctx context.Context, id string) (*domain.HealthRecord, error)
	Update(ctx context.Context, id string, update_input servicedto.HealthRecordUpdateInput) error
	ListByUserID(ctx context.Context, userId string) ([]*domain.HealthRecord, error)
	Delete(ctx context.Context, id string) error
}

type HealthRecordHandler struct {
	healthRecordService HealthRecordService
}

func NewHealthRecordHandler(service HealthRecordService) *HealthRecordHandler {
	return &HealthRecordHandler{healthRecordService: service}
}

// Create godoc
//
// @Summary      Create Health Record
// @Description  Creates a new Health Record
// @Tags         Health Record
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateHealthRecordRequest true "Health Record Type data"
// @Success      201 {object} dto.HealthRecordIdResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /api/healthrecord [post]
func (h *HealthRecordHandler) Create(c *gin.Context) {
	var req handlerdto.CreateHealthRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, handlerdto.ErrorResponse{Error: err.Error()})
		return
	}

	id, err := h.healthRecordService.Create(c.Request.Context(), toHealthRecordServiceCreateInput(req))

	if err != nil {
		if errors.Is(err, domain.ErrInvalidUserOrHealthRecordType) {
			c.JSON(http.StatusBadRequest, handlerdto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, handlerdto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, handlerdto.HealthRecordIdResponse{Id: id})

}

// Update godoc
//
// @Summary      Update Health Record
// @Description  Updates a Health Record
// @Tags         Health Record
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "Health Record UUID"
// @Param        request body dto.UpdateHealthRecordRequest true "Health Record Type data"
// @Success      201 {object} dto.HealthRecordIdResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /api/healthrecord/{id} [patch]
func (h *HealthRecordHandler) Update(c *gin.Context) {
	var req handlerdto.UpdateHealthRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, handlerdto.ErrorResponse{Error: err.Error()})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, handlerdto.ErrorResponse{Error: "missing record id"})
		return
	}

	err := h.healthRecordService.Update(c.Request.Context(), id, toHealthRecordServiceUpdateInput(req))
	if err != nil {
		if errors.Is(err, domain.ErrHealthRecordNotFound) {
			c.JSON(http.StatusNotFound, handlerdto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, handlerdto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, handlerdto.HealthRecordIdResponse{Id: id})

}

// GetByID godoc
//
// @Summary      Get Health Record By ID
// @Description  Gets a Health Record by its ID
// @Tags         Health Record
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "Health Record UUID"
// @Success      200 {object} dto.HealthRecord
// @Failure      400 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /api/healthrecord/{id} [get]
func (h *HealthRecordHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, handlerdto.ErrorResponse{Error: "missing record id"})
		return
	}

	healthRecord, err := h.healthRecordService.GetByID(c.Request.Context(), id)

	if err != nil {
		if errors.Is(err, domain.ErrHealthRecordNotFound) {
			c.JSON(http.StatusNotFound, handlerdto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, handlerdto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, toHealthRecordResponse(healthRecord))

}

// ListByUserID godoc
//
// @Summary      List Health Records by User ID
// @Description  Lists Health Records by User ID
// @Tags         Health Record
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        user_id query string true "User UUID"
// @Success      200 {object} []dto.HealthRecord
// @Failure      400 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /api/healthrecord/list [get]
func (h *HealthRecordHandler) ListByUserID(c *gin.Context) {
	userID := c.Query("user_id")

	if userID == "" {
		c.JSON(http.StatusBadRequest, handlerdto.ErrorResponse{Error: domain.ErrorUserIDNotProvided.Error()})
		return
	}

	healthRecords, err := h.healthRecordService.ListByUserID(c.Request.Context(), userID)

	if err != nil {
		if errors.Is(err, domain.ErrHealthRecordNotFound) {
			c.JSON(http.StatusNotFound, handlerdto.ErrorResponse{Error: err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, handlerdto.ErrorResponse{Error: err.Error()})
		return
	}

	if len(healthRecords) == 0 {
		c.JSON(http.StatusNotFound, handlerdto.ErrorResponse{Error: domain.ErrHealthRecordNotFound.Error()})
		return
	}
	var healthRecordList []handlerdto.HealthRecord

	for _, hr := range healthRecords {
		healthRecordList = append(healthRecordList, toHealthRecordResponse(hr))
	}
	c.JSON(http.StatusOK, healthRecordList)
}

// Delete godoc
//
// @Summary      Delete Health Record By ID
// @Description  Deletes a Health Record by its ID
// @Tags         Health Record
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id query string true "Health Record UUID"
// @Success      204
// @Failure      400 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /api/healthrecord [delete]
func (h *HealthRecordHandler) Delete(c *gin.Context) {
	id := c.Query("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, handlerdto.ErrorResponse{Error: domain.ErrEmptyId.Error()})
		return
	}

	err := h.healthRecordService.Delete(c.Request.Context(), id)

	if err != nil {
		if errors.Is(err, domain.ErrHealthRecordNotFound) {
			c.JSON(http.StatusNotFound, handlerdto.ErrorResponse{Error: err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, handlerdto.ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func toHealthRecordServiceCreateInput(req handlerdto.CreateHealthRecordRequest) servicedto.HealthRecordCreateInput {
	return servicedto.HealthRecordCreateInput{
		UserID:             req.UserID,
		HealthRecordTypeID: req.HealthRecordTypeID,
		Value:              req.Value,
		Notes:              req.Notes,
		RecordedAt:         req.RecordedAt,
	}
}

func toHealthRecordServiceUpdateInput(req handlerdto.UpdateHealthRecordRequest) servicedto.HealthRecordUpdateInput {
	return servicedto.HealthRecordUpdateInput{
		Value:      req.Value,
		Notes:      req.Notes,
		RecordedAt: req.RecordedAt,
	}
}

func toHealthRecordResponse(resp *domain.HealthRecord) handlerdto.HealthRecord {
	return handlerdto.HealthRecord{
		ID:                 resp.ID,
		UserID:             resp.UserID,
		HealthRecordTypeID: resp.HealthRecordTypeID,
		Value:              resp.Value,
		RecordedAt:         resp.RecordedAt,
		Notes:              resp.Notes,
		CreatedAt:          resp.CreatedAt,
		UpdatedAt:          resp.UpdatedAt,
	}
}
