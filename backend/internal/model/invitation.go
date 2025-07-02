package model

import "github.com/google/uuid"

type Invitation struct {
	Base
	Token  []byte    `gorm:"uniqueIndex;not null" json:"token"`
	UserID uuid.UUID `gorm:"not null" json:"userId"`
	User   User      `gorm:"foreignKey:UserID" json:"user"`
}
