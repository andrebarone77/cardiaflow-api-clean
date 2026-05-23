package domain

import "time"

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
