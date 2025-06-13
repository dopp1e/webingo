package model

import "github.com/google/uuid"

type Game struct {
	Base
	YoutubeVideoID *string `json:"youtubeVideoId,omitempty"`
	Secret         *string `json:"-"` // Typically hidden from JSON output
	Height         int     `json:"height"`
	Width          int     `json:"width"`
	Finished       bool    `json:"finished"`
	Privacy        *bool   `json:"private"`
	Password       *string `json:"-"` // Password for private games, hidden from JSON output

	BoardID uuid.UUID `json:"boardId"`
	Board   Board     `gorm:"foreignKey:BoardID"`
	UserID  uuid.UUID `json:"userId"` // User who created the game
	User    User      `gorm:"foreignKey:UserID"`

	Cells       []Cell        `gorm:"foreignKey:GameID" json:"cells"`
	PlayerMoves []PlayerMove  `gorm:"foreignKey:GameID" json:"playerMoves"`
	Comments    []GameComment `gorm:"foreignKey:GameID" json:"comments"`
	Votes       []GameVote    `gorm:"foreignKey:GameID" json:"votes"`
}
