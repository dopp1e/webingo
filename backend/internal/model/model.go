package model

import (
	"gorm.io/gorm"
)

// context key strings
type boardkey string
type contextModel struct {
	Board boardkey
}

var Context = contextModel{
	Board: "board",
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
	)

	if err != nil {
		return err
	}

	return nil
}
