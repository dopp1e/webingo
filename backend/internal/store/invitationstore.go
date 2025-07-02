package store

import (
	"context"

	"github.com/dopp1e/webingo/backend/internal/model"
	"gorm.io/gorm"
)

type InvitationStore struct {
	db *gorm.DB
}

func (s *InvitationStore) Create(ctx context.Context, invitation *model.Invitation) error {
	result := s.db.WithContext(ctx).Create(invitation)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
