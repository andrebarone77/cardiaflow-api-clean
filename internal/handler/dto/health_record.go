package dto

import "time"

type CreateHealthRecordRequest struct {
	HealthRecordTypeID *string   `json:"health_record_type_id" binding:"required"`
	Value              *float64  `json:"value" binding:"required"`
	Notes              string    `json:"notes"`
	UserID             *string   `json:"user_id" binding:"required"`
	RecordedAt         time.Time `json:"recorded_at"`
}

type UpdateHealthRecordRequest struct {
	Value      *float64   `json:"value"`
	Notes      *string    `json:"notes"`
	RecordedAt *time.Time `json:"recorded_at"`
}
