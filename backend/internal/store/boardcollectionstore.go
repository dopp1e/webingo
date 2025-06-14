package store

import (
	"context"

	"github.com/dopp1e/webingo/backend/internal/model"
	"gorm.io/gorm"
)

type BoardCollectionStore struct {
	db *gorm.DB
}

func (s *BoardCollectionStore) Create(ctx context.Context, collection *model.BoardCollection) error {
	result := s.db.WithContext(ctx).Create(collection)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
