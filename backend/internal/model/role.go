package model

type Role struct {
	Base
	Name        string `gorm:"uniqueIndex;not null" json:"name"`
	Description string `gorm:"not null" json:"description"`
	Level       int    `gorm:"not null" json:"level"`
	Users       []User `gorm:"foreignKey:RoleID"`
}
