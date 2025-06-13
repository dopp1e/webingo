package model

import "github.com/google/uuid"

// BoardComment represents a comment specifically on a board.
type BoardComment struct {
	Base
	Content string    `gorm:"not null" json:"content"` // Content of the comment
	BoardID uuid.UUID `json:"boardId"`
	Board   Board     `gorm:"foreignKey:BoardID"`
	UserID  uuid.UUID `json:"userId"`
	User    User      `gorm:"foreignKey:UserID"`
}

// GameComment represents a comment specifically on a game.
type GameComment struct {
	Base
	Content string    `gorm:"not null" json:"content"` // Content of the comment
	GameID  uuid.UUID `json:"gameId"`
	Game    Game      `gorm:"foreignKey:GameID"`
	UserID  uuid.UUID `json:"userId"`
	User    User      `gorm:"foreignKey:UserID"`
}
