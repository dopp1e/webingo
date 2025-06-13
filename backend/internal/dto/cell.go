package dto

import "github.com/google/uuid"

// CellCreateRequest for creating a cell (e.g., when setting up a game).
type CellCreateRequest struct {
	SpaceID uuid.UUID `json:"spaceId" validate:"required"`
	PosH    int       `json:"posH" validate:"required,min=0,max=100"`
	PosW    int       `json:"posW" validate:"required,min=0,max=100"`
}
