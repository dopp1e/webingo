package store

import (
	"context"

	"github.com/dopp1e/webingo/backend/internal/model"
	"github.com/google/uuid"
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

func (s *BoardCommentStore) GetByID(ctx context.Context, id uuid.UUID) (*model.BoardComment, error) {
	var comment model.BoardComment
	result := s.db.WithContext(ctx).First(&comment, "id = ?", id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &comment, nil
}
