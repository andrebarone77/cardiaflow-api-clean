package handler

import (
	"errors"
	"net/http"

	"github.com/andrebarone77/cardiaflow-api/internal/domain"
	handlerdto "github.com/andrebarone77/cardiaflow-api/internal/handler/dto"
	"github.com/andrebarone77/cardiaflow-api/internal/service"
	servicedto "github.com/andrebarone77/cardiaflow-api/internal/service/dto"
	"github.com/gin-gonic/gin"
)

type HealthRecordHandler struct {
	healthRecordService *service.HealthRecordService
}

func NewHealthRecordHandler(service *service.HealthRecordService) *HealthRecordHandler {
	return &HealthRecordHandler{healthRecordService: service}
}

func (h *HealthRecordHandler) Create(c *gin.Context) {
	var req handlerdto.CreateHealthRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.healthRecordService.Create(c.Request.Context(), toHealthRecordServiceCreateInput(req))

	if err != nil {
		if errors.Is(err, domain.ErrInvalidUserOrHealthRecordType) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})

}

func (h *HealthRecordHandler) Update(c *gin.Context) {
	var req handlerdto.UpdateHealthRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing record id"})
		return
	}

	err := h.healthRecordService.Update(c.Request.Context(), id, toHealthRecordServiceUpdateInput(req))
	if err != nil {
		if errors.Is(err, domain.ErrHealthRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})

}

func (h *HealthRecordHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	healthRecord, err := h.healthRecordService.GetByID(c.Request.Context(), id)

	if err != nil {
		if errors.Is(err, domain.ErrHealthRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, healthRecord)

}

func (h *HealthRecordHandler) ListByUserID(c *gin.Context) {
	userID := c.Query("user_id")

	if userID == "" {
		c.JSON(http.StatusOK, gin.H{"error": domain.ErrorUserIDNotProvided})
		return
	}

	healthRecords, err := h.healthRecordService.ListByUserID(c.Request.Context(), userID)

	if err != nil {
		if errors.Is(err, domain.ErrHealthRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(healthRecords) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": domain.ErrHealthRecordNotFound.Error()})
		return
	}

	c.JSON(http.StatusOK, healthRecords)
}

func (h *HealthRecordHandler) Delete(c *gin.Context) {
	id := c.Query("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrEmptyId.Error()})
		return
	}

	err := h.healthRecordService.Delete(c.Request.Context(), id)

	if err != nil {
		if errors.Is(err, domain.ErrHealthRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
