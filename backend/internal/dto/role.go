package dto

type RoleCreateRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Level       int    `json:"level" validate:"required"`
}
