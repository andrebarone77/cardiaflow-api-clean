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

type HealthRecordIdResponse struct {
	Id string `json: id`
}

type HealthRecord struct {
	ID                 string    `json:"id"`
	UserID             string    `json:"user_id"`
	HealthRecordTypeID string    `json:"health_record_type_id"`
	Value              float64   `json:"value"`
	RecordedAt         time.Time `json:"recorded_at"`
	Notes              *string   `json:"notes,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
