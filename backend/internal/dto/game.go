package dto

import "github.com/google/uuid"

// GameCreateRequest for creating a new game.
type GameCreateRequest struct {
	YoutubeVideoID *string    `json:"youtubeVideoId,omitempty"`
	Secret         *string    `json:"secret" validate:"required"`
	Height         int        `json:"height" validate:"required,min=1,max=100"`
	Width          int        `json:"width" validate:"required,min=1,max=100"`
	Finished       bool       `json:"finished"`
	Privacy        *bool      `json:"private" validate:"required"` // Must be provided
	Password       *string    `json:"password,omitempty"`
	BoardID        uuid.UUID  `json:"boardId" validate:"required"` // The board the game is associated with
	UserID         *uuid.UUID `json:"userId,omitempty"`            // NOW A POINTER: Optional for anonymous games
	// Cells and PlayerMoves are typically created/managed via separate endpoints/requests after game creation.
}

// GameUpdateRequest for updating an existing game.
type GameUpdateRequest struct {
	YoutubeVideoID *string    `json:"youtubeVideoId,omitempty"`
	Secret         *string    `json:"secret" validate:"required"` // Still required for update? Or omit?
	Height         int        `json:"height" validate:"required,min=1,max=100"`
	Width          int        `json:"width" validate:"required,min=1,max=100"`
	Finished       bool       `json:"finished"`
	Privacy        *bool      `json:"private" validate:"required"`
	Password       *string    `json:"password,omitempty"`
	Version        int        `json:"version" validate:"required"` // For optimistic locking
	UserID         *uuid.UUID `json:"userId,omitempty"`            // NOW A POINTER: Allow changing owner if permitted
}
