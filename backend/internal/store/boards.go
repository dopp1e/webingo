package store

import (
	"context"

	"github.com/dopp1e/webingo/backend/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BoardStore struct {
	db *gorm.DB
}

func (s *BoardStore) Create(ctx context.Context, board *model.Board) error {
	result := s.db.WithContext(ctx).Create(board)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *BoardStore) GetById(ctx context.Context, boardId uuid.UUID) (*model.Board, error) {
	var board model.Board
	result := s.db.WithContext(ctx).Preload("Spaces").First(&board, "id = ?", boardId)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &board, nil
}
