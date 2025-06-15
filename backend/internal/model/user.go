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

	BoardCollections []BoardCollection `gorm:"foreignKey:UserID"`
	Notifications    []Notification    `gorm:"foreignKey:UserID"`
	Reports          []Report          `gorm:"foreignKey:ReporterUserID"` // Reports made by this user

	Following []User `gorm:"many2many:follows;joinForeignKey:FollowerID;joinReferences:FollowedUserID"`
	Followers []User `gorm:"many2many:follows;joinForeignKey:FollowedUserID;joinReferences:FollowerID"`
}
