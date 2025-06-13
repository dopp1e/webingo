package model

import "github.com/google/uuid"

// PlayerMove represents a move made by a player in the game.
type PlayerMove struct {
	Base
	CellH     int       `json:"cellH"`
	CellW     int       `json:"cellW"`
	Timestamp int       `json:"timestamp"`
	GameID    uuid.UUID `json:"gameId"`
	Game      Game      `gorm:"foreignKey:GameID"`
}
