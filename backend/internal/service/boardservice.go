package service

import (
	"context"

	"github.com/dopp1e/webingo/backend/internal/model"
	"github.com/dopp1e/webingo/backend/internal/store"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BoardService struct {
	storage store.Storage
	db      *gorm.DB
}

func (s *BoardService) CreateBoard(ctx context.Context, board *model.Board) error {
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback() // defer rollback in case of error

	var preparedSpaces []model.Space
	for _, incomingSpace := range board.Spaces {
		foundOrCreatedSpace, err := s.storage.Spaces.FindOrCreateByContent(ctx, &incomingSpace)
		if err != nil {
			return err
		}
		preparedSpaces = append(preparedSpaces, *foundOrCreatedSpace)
	}

	board.Spaces = preparedSpaces

	if err := s.storage.Boards.Create(ctx, board); err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (s *BoardService) GetById(ctx context.Context, boardId uuid.UUID) (*model.Board, error) {
	board, err := s.storage.Boards.GetById(ctx, boardId)
	if err != nil {
		return nil, err
	}

	return board, nil
}
