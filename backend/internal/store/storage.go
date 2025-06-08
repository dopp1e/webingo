package store

import (
	"context"

	"github.com/dopp1e/webingo/backend/internal/model"
	"gorm.io/gorm"
)

type Storage struct {
	Boards interface {
		Create(context.Context, *model.Board) error
	}
	Users interface {
		Create(context.Context, *model.User) error
		GetByUsername(context.Context, string) (*model.User, error)
		Exists(context.Context, string) (bool, error)
	}
}

func NewStorage(db *gorm.DB) Storage {
	return Storage{
		Users:  &UsersStore{db},
		Boards: &BoardStore{db},
	}
}
