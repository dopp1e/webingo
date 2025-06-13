package model

import "github.com/google/uuid"

type Cell struct {
	Base
	SpaceID uuid.UUID `json:"spaceId"`
	Space   Space     `gorm:"foreignKey:SpaceID"`
	PosH    int       `json:"posH"`
	PosW    int       `json:"posW"`
	GameID  uuid.UUID `json:"gameId"`
	Game    Game      `gorm:"foreignKey:GameID"`
}
