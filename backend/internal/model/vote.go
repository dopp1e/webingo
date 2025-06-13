package model

import "github.com/google/uuid"

// BoardVote represents an upvote/downvote on a Board.
type BoardVote struct {
	Base
	BoardID  uuid.UUID `gorm:"not null;uniqueIndex:idx_board_user_vote,composite:BoardID,UserID" json:"boardId"`
	Board    Board     `gorm:"foreignKey:BoardID"`
	UserID   uuid.UUID `gorm:"not null" json:"userId"`
	User     User      `gorm:"foreignKey:UserID"`
	IsUpvote bool      `gorm:"not null" json:"isUpvote"`
}

// GameVote represents an upvote/downvote on a Game.
type GameVote struct {
	Base
	GameID   uuid.UUID `gorm:"not null;uniqueIndex:idx_game_user_vote,composite:GameID,UserID" json:"gameId"`
	Game     Game      `gorm:"foreignKey:GameID"`
	UserID   uuid.UUID `gorm:"not null" json:"userId"`
	User     User      `gorm:"foreignKey:UserID"`
	IsUpvote bool      `gorm:"not null" json:"isUpvote"`
}
