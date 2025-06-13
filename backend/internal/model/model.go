package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base struct for common fields across models.
type Base struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `gorm:"index" json:"-"` // Soft delete field
}

// BoardPayload represents the data that is accepted when creating a board.
type BoardPayload struct {
	Name        string    `json:"name" validate:"required,min=3,max=100"`
	Description string    `json:"description" validate:"max=500"`
	Version     int       `gorm:"default:1;column:version" json:"version"`               // Version for optimistic locking
	UserID      uuid.UUID `json:"userId"`                                                // Foreign key to the User who created the board
	User        User      `gorm:"foreignKey:UserID" json:"user"`                         // User who created the board
	Tags        []string  `gorm:"type:text[]" json:"tags" validate:"dive,min=1,max=20"`  // Tags for the board
	Spaces      []Space   `gorm:"many2many:board_spaces;" json:"spaces" validate:"dive"` // Spaces on the board
}

// Board represents a bingo board.
type Board struct {
	Base
	BoardPayload
}

// Space represents a space on the bingo board.
type Space struct {
	Base
	Content string  `gorm:"unique;not null" json:"content" validate:"required,min=1,max=100"`
	Boards  []Board `gorm:"many2many:board_spaces;" json:"boards"`
}

type Comment struct {
	Base
	BoardID uuid.UUID `json:"boardId"` // Foreign key to the Board
	Board   Board     `gorm:"foreignKey:BoardID"`
	Content string    `gorm:"not null" json:"content" validate:"required,min=1,max=500"` // Content of the comment
	UserID  uuid.UUID `json:"userId"`                                                    // Foreign key to the User who made the comment
	User    User      `gorm:"foreignKey:UserID"`                                         // User who made the comment
}

// Game represents a single player's bingo game instance.
type Game struct {
	Base
	BoardID        uuid.UUID    `json:"board_id"` // Foreign key to the Board
	Board          Board        `gorm:"foreignKey:BoardID"`
	YoutubeVideoID *string      `json:"youtubeVideoId,omitempty"` // Optional YouTube video ID for the game
	Secret         string       `json:"secret"`                   // Secret for the game, used to disallow updating the game by other players
	Height         int          `json:"height"`                   // Height of the bingo board
	Width          int          `json:"width"`                    // Width of the bingo board
	Finished       bool         `json:"finished"`                 // Indicates if the game is finished
	Privacy        bool         `json:"private"`                  // Indicates if the game is private
	Password       *string      `json:"password,omitempty"`       // Optional password for private games
	Cells          []Cell       `gorm:"foreignKey:GameID"`        // Cells in the game
	PlayerMoves    []PlayerMove `gorm:"foreignKey:GameID"`
}

// Cell represents a cell in the bingo game, linking a space to a specific position on the board.
type Cell struct {
	Base
	SpaceID uuid.UUID `json:"spaceId"` // Foreign key to the Space
	Space   Space     `gorm:"foreignKey:SpaceID"`
	GameID  uuid.UUID `json:"gameId"` // Foreign key to the Game
	Game    Game      `gorm:"foreignKey:GameID"`
	PosH    int       `json:"posH"` // Horizontal position on the board
	PosW    int       `json:"posW"` // Vertical position on the board
}

// PlayerMove represents a move made by a player in the game.
type PlayerMove struct {
	Base
	// Foreign key to the Game instance this move belongs to.
	GameID    uuid.UUID `json:"gameId"`
	Game      Game      `gorm:"foreignKey:GameID"` // Belongs To relationship (GameID in PlayerMove table)
	CellH     int       `json:"cellH"`             // Horizontal position of the cell marked
	CellW     int       `json:"cellW"`             // Vertical position of the cell marked
	Timestamp int       `json:"timestamp"`         // Timestamp of the move (e.g., seconds into the video)
}

type Role struct {
	Base
	Name        string `gorm:"uniqueIndex;not null" json:"name"` // Unique name for the role
	Description string `gorm:"not null" json:"description"`      // Description of the role
	Level       int    `gorm:"not null" json:"level"`            // Level of the role, used for permissions
}

// User represents a user in the system.
type User struct {
	Base
	Username string    `gorm:"uniqueIndex;not null" json:"username"` // Unique username for the user
	Email    string    `gorm:"uniqueIndex;not null" json:"email"`    // Unique email for the user
	Password string    `gorm:"not null" json:"-"`                    // Password for the user, not exposed in JSON
	IsActive bool      `gorm:"default:true" json:"isActive"`         // Indicates if the user is active
	RoleID   uuid.UUID `json:"roleId"`                               // Foreign key to the Role
	Role     Role      `gorm:"foreignKey:RoleID"`                    // Role of the user
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
