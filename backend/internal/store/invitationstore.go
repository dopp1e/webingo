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

func (s *InvitationStore) GetByToken(ctx context.Context, token []byte) (*model.Invitation, error) {
	var invitation model.Invitation
	result := s.db.WithContext(ctx).Where("token = ?", token).First(&invitation)
	if result.Error != nil {
		return nil, result.Error
	}

	return &invitation, nil
}

func (s *InvitationStore) Delete(ctx context.Context, invitation *model.Invitation) error {
	result := s.db.WithContext(ctx).Delete(invitation)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
