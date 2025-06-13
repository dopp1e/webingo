package dto

// SpaceRequest represents the input for a single Space (e.g., within a BoardCreateRequest).
type SpaceRequest struct {
	Content string `json:"content" validate:"required,min=1,max=100"`
}
