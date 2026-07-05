package dto

type ErrorResponse struct {
	Error string `json:"error" example:"User not found"`
}
