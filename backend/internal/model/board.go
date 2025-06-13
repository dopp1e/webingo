package model

import "github.com/google/uuid"

// Board represents a bingo board with its associated metadata and relationships.
type Board struct {
	Base
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Version     int      `gorm:"default:1;column:version" json:"version"` // For optimistic locking
	Tags        []string `gorm:"type:text[]" json:"tags"`                 // Tags for the board

	UserID   uuid.UUID      `json:"userId"`
	User     User           `gorm:"foreignKey:UserID"`
	Comments []BoardComment `gorm:"foreignKey:BoardID"`
	Votes    []BoardVote    `gorm:"foreignKey:BoardID"`
	Spaces   []Space        `gorm:"many2many:board_spaces;"` // GORM tag, no json tag if Space has Board reference
}
