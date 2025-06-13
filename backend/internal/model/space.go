package model

// Space represents a space on the bingo board.
type Space struct {
	Base
	Content string  `gorm:"unique;not null" json:"content"`
	Boards  []Board `gorm:"many2many:board_spaces;" json:"-"` // json:"-" to prevent JSON recursion
}
