package dto

type CreateHealthRecordTypeRequest struct {
	Name string  `json:"name" binding:"required"`
	Code string  `json:"code" binding:"required"`
	Unit *string `json:"unit"`
}

type UpdateHealthRecordTypeRequest struct {
	Name *string `json:"name"`
	Code *string `json:"code"`
	Unit *string `json:"unit"`
}
