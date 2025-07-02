package model

import (
	"gorm.io/gorm"
)

// context key strings
type boardkey string
type userkey string
type contextModel struct {
	Board boardkey
	User  userkey
}

var Context = contextModel{
	Board: "board",
	User:  "user",
}

// Migrate performs the database migrations for the bingo application models.
func Migrate(db *gorm.DB) error {
	// Automatically migrate the schema, creating tables and relationships.
	err := db.AutoMigrate(
		&Board{},
		&Space{},
		&Game{},
		&Cell{},
		&PlayerMove{},
		&User{},
		&BoardComment{},
		&GameComment{},
		&Role{},
		&BoardVote{},
		&GameVote{},
		&BoardCollection{},
		&Notification{},
		&Report{},
		&Follow{},
		&Activity{},
		&Invitation{},
	)

	if err != nil {
		return err
	}

	return nil
}
