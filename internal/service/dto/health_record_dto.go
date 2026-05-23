package dto

import "time"

type HealthRecordCreateInput struct {
	UserID             *string
	HealthRecordTypeID *string
	Value              *float64
	Notes              string
	RecordedAt         time.Time
}

type HealthRecordUpdateInput struct {
	Value      *float64
	Notes      *string
	RecordedAt *time.Time
}
