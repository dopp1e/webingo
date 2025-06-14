package model

import "github.com/google/uuid"

type Cell struct {
	Base
	SpaceID uuid.UUID `json:"spaceId"`
	Space   Space     `gorm:"foreignKey:SpaceID" json:"-"` // Hide Space in JSON by default
	PosH    int       `json:"posH"`
	PosW    int       `json:"posW"`
	GameID  uuid.UUID `json:"gameId"`
	Game    Game      `gorm:"foreignKey:GameID" json:"-"` // Hide Game in JSON by default
}
