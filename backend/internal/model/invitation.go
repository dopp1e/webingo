package model

import (
	"time"

	"github.com/google/uuid"
)

type Invitation struct {
	Base
	Token  []byte    `gorm:"uniqueIndex;not null" json:"token"`
	UserID uuid.UUID `gorm:"not null" json:"userId"`
	Expiry time.Time `gorm:"not null" json:"expiry"`
	User   User      `gorm:"foreignKey:UserID" json:"user"`
}
