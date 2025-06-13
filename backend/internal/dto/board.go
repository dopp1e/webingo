package dto

// BoardCreateRequest for creating a new board.
type BoardCreateRequest struct {
	Name        string         `json:"name" validate:"required,min=3,max=100"`
	Description string         `json:"description" validate:"max=500"`
	Tags        []string       `json:"tags" validate:"dive,min=1,max=20"`
	Spaces      []SpaceRequest `json:"spaces" validate:"dive"` // Spaces input as SpaceRequest DTOs
}

// BoardUpdateRequest for updating an existing board (PUT).
// Version for optimistic locking will be part of the request, but not in the payload.
type BoardUpdateRequest struct {
	Name        string         `json:"name" validate:"required,min=3,max=100"`
	Description string         `json:"description" validate:"max=500"`
	Version     int            `json:"version" validate:"required"` // Version must be provided for optimistic locking
	Tags        []string       `json:"tags" validate:"dive,min=1,max=20"`
	Spaces      []SpaceRequest `json:"spaces" validate:"dive"` // New set of spaces for PUT
}
