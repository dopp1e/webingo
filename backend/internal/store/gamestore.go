package store

import (
	"context"

	"github.com/dopp1e/webingo/backend/internal/model"
	"gorm.io/gorm"
)

type GameStore struct {
	db *gorm.DB
}

func (s *GameStore) Create(ctx context.Context, game *model.Game) error {
	result := s.db.WithContext(ctx).Create(game)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
