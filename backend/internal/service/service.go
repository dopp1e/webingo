package service

import (
	"context"

	"github.com/dopp1e/webingo/backend/internal/model"
	"github.com/dopp1e/webingo/backend/internal/store"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	storage store.Storage
	db      *gorm.DB
	Boards  interface {
		CreateBoard(context.Context, *model.Board) error
		GetById(context.Context, uuid.UUID) (*model.Board, error)
	}
	Users interface {
		Create(context.Context, *model.User) error
		GetByUsername(context.Context, string) (*model.User, error)
		Exists(context.Context, string) (bool, error)
		CreateIfNotExists(context.Context, *model.User) (*model.User, error)
	}
}

func NewService(storage store.Storage, db *gorm.DB) *Service {
	return &Service{
		storage: storage,
		db:      db,
		Boards: &BoardService{
			storage: storage,
			db:      db,
		},
		Users: &UserService{
			storage: storage,
			db:      db,
		},
	}
}
