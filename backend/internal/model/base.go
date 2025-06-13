package model

import (
	"time"

	"github.com/google/uuid"
)

// Base struct for common fields across models.
type Base struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `gorm:"index" json:"-"` // Soft delete
}
