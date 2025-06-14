package model

import "github.com/google/uuid"

// CollectionType defines the type of a collection (e.g., "FAVORITES", "USER_CREATED")
type BoardCollectionType string

const (
	BoardCollectionTypeFavorites   BoardCollectionType = "FAVORITES"
	BoardCollectionTypeUserCreated BoardCollectionType = "USER_CREATED"
)

// BoardCollection represents a user-curated collection or playlist of boards.
type BoardCollection struct {
	Base
	Name        string              `gorm:"not null" json:"name"`          // Name of the collection (e.g., "My Party Boards")
	Description string              `json:"description,omitempty"`         // Optional description
	UserID      uuid.UUID           `gorm:"not null;index" json:"userId"`  // User who owns the collection
	User        User                `gorm:"foreignKey:UserID" json:"-"`    // Hide User in JSON
	Type        BoardCollectionType `gorm:"not null" json:"type"`          // e.g., "FAVORITES", "USER_CREATED"
	IsPublic    bool                `gorm:"default:false" json:"isPublic"` // Can others view this collection?

	// Many-to-many relationship with Boards
	Boards []Board `gorm:"many2many:collection_boards;" json:"boards"`
}
