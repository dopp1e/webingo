package model

import (
	"gorm.io/gorm"
)

// BoardPayload represents the data that is accepted when creating a board.
type BoardPayload struct {
	Name        string   `json:"name" validate:"required,min=3,max=100"`
	Description string   `json:"description" validate:"max=500"`
	Version     int      `gorm:"default:1;column:version" json:"version"`               // Version for optimistic locking
	Tags        []string `gorm:"type:text[]" json:"tags" validate:"dive,min=1,max=20"`  // Tags for the board
	Spaces      []Space  `gorm:"many2many:board_spaces;" json:"spaces" validate:"dive"` // Spaces on the board
}

type PlayerMovePayload struct {
	CellH     int `json:"cellH"`     // Horizontal position of the cell marked
	CellW     int `json:"cellW"`     // Vertical position of the cell marked
	Timestamp int `json:"timestamp"` // Timestamp of the move (e.g., seconds into the video)
}

// UserPayloadIn represents the data that is accepted when creating a new user.
type UserPayloadIn struct {
	Username string `json:"username" validate:"required,min=3,max=50"` // Username for the user
	Email    string `json:"email" validate:"required,email"`           // Email for the user
	Password string `json:"password" validate:"required"`              // Password for the user
}

// UserPayload represents the data that is accepted when sending user information.
type UserPayload struct {
	Username string `json:"username" validate:"required,min=3,max=50"` // Username for the user
	Email    string `json:"email" validate:"required,email"`           // Email for the user
	Password string `json:"-" validate:"omitempty"`                    // Password for the user, optional for updates
}

// context key strings
type boardkey string
type contextModel struct {
	Board boardkey
}

var Context = contextModel{
	Board: "board",
}

// Migrate performs the database migrations for the bingo application models.
func Migrate(db *gorm.DB) error {
	// Automatically migrate the schema, creating tables and relationships.
	err := db.AutoMigrate(&Board{}, &Space{}, &Game{}, &Cell{}, &PlayerMove{}, &User{})

	if err != nil {
		return err
	}

	return nil
}
