package service

import (
	"context"

	"github.com/dopp1e/webingo/backend/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	db     *gorm.DB
	Boards interface {
		CreateBoard(context.Context, *model.Board) error
		GetById(context.Context, uuid.UUID) (*model.Board, error)
		UpdateBoard(context.Context, *model.Board) error
		DeleteBoard(context.Context, uuid.UUID) error
	}
	Users interface {
		Create(context.Context, *model.User) error
		GetByUsername(context.Context, string) (*model.User, error)
		Exists(context.Context, string) (bool, error)
		CreateIfNotExists(context.Context, *model.User) (*model.User, error)
	}
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		db: db,
		Boards: &BoardService{
			db: db,
		},
		Users: &UserService{
			db: db,
		},
	}
}
