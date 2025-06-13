package dto

import "github.com/google/uuid"

// BoardCollectionCreateRequest for creating a new board collection.
type BoardCollectionCreateRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Description string `json:"description,omitempty" validate:"max=500"`
	IsPublic    bool   `json:"isPublic"`
	// Type is usually set by the server (e.g., "FAVORITES" or "USER_CREATED")
}

// BoardCollectionUpdateRequest for updating an existing board collection.
type BoardCollectionUpdateRequest struct {
	Name        string `json:"name" validate:"omitempty,min=3,max=100"` // Omitempty for partial updates
	Description string `json:"description,omitempty" validate:"max=500"`
	IsPublic    *bool  `json:"isPublic,omitempty"` // Use pointer for partial update
	// Type cannot be changed via update
	Version int `json:"version" validate:"required"` // For optimistic locking
}

// BoardCollectionAddBoardRequest for adding a board to a collection.
type BoardCollectionAddBoardRequest struct {
	BoardID uuid.UUID `json:"boardId" validate:"required"`
}

// BoardCollectionRemoveBoardRequest for removing a board from a collection.
type BoardCollectionRemoveBoardRequest struct {
	BoardID uuid.UUID `json:"boardId" validate:"required"`
}
