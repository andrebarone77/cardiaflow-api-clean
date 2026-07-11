package dto

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UpdateUserRequest struct {
	Name     *string `json:"name"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
	ID       *string `json:"id"`
}

type UserResponse struct {
	ID    string `json:"id" example:"3dcb50de-a7d4-4d1d-8e78-6c8f1d0e7c4b"`
	Name  string `json:"name" example:"Andre Barone"`
	Email string `json:"email" example:"andre@email.com"`
}
