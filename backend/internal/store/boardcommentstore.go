package store

import (
	"context"

	"github.com/dopp1e/webingo/backend/internal/model"
	"gorm.io/gorm"
)

type BoardCommentStore struct {
	db *gorm.DB
}

func (s *BoardCommentStore) Create(ctx context.Context, comment *model.BoardComment) error {
	result := s.db.WithContext(ctx).Create(comment)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
