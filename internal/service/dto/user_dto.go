package dto

type CreateUserInput struct {
	Name     string
	Email    string
	Password string
}

type UpdateUserInput struct {
	Name     *string
	Email    *string
	Password *string
}
