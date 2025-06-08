package store

import (
	"context"

	"github.com/dopp1e/webingo/backend/internal/model"
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
