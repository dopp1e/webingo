package model

import "github.com/google/uuid"

type User struct {
	Base
	Username string `gorm:"uniqueIndex;not null" json:"username"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `gorm:"not null" json:"-"` // Password stored here, not exposed in JSON output
	IsActive bool   `gorm:"default:true" json:"isActive"`

	RoleID uuid.UUID `json:"roleId"`
	Role   Role      `gorm:"foreignKey:RoleID"`

	Boards        []Board        `gorm:"foreignKey:UserID"`
	Games         []Game         `gorm:"foreignKey:UserID"`
	BoardComments []BoardComment `gorm:"foreignKey:UserID"`
	GameComments  []GameComment  `gorm:"foreignKey:UserID"`
	BoardVotes    []BoardVote    `gorm:"foreignKey:UserID"`
	GameVotes     []GameVote     `gorm:"foreignKey:UserID"`
}
