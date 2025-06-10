package store

import (
	"context"

	"github.com/dopp1e/webingo/backend/internal/model"
	"gorm.io/gorm"
)

type SpaceStore struct {
	db *gorm.DB
}

func (s *SpaceStore) FindOrCreateByContent(ctx context.Context, space *model.Space) (*model.Space, error) {
	foundSpace := model.Space{}

	result := s.db.WithContext(ctx).Where(model.Space{Content: space.Content}).FirstOrCreate(&foundSpace, space)
	if result.Error != nil {
		return nil, result.Error
	}

	return &foundSpace, nil // New space created
}
