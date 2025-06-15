package store

import (
	"context"

	"github.com/dopp1e/webingo/backend/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FollowStore struct {
	db *gorm.DB
}

func (s *FollowStore) Create(ctx context.Context, follow *model.Follow) error {
	result := s.db.WithContext(ctx).Create(follow)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *FollowStore) GetFollowers(ctx context.Context, userID uuid.UUID) ([]model.Follow, error) {
	var follows []model.Follow
	result := s.db.WithContext(ctx).Where("followed_id = ?", userID).Find(&follows)

	if result.Error != nil {
		return nil, result.Error
	}

	return follows, nil
}

func (s *FollowStore) GetFollowing(ctx context.Context, userID uuid.UUID) ([]model.Follow, error) {
	var follows []model.Follow
	result := s.db.WithContext(ctx).Where("follower_id = ?", userID).Find(&follows)

	if result.Error != nil {
		return nil, result.Error
	}

	return follows, nil
}

func (s *FollowStore) Delete(ctx context.Context, follow *model.Follow) error {
	result := s.db.WithContext(ctx).Delete(follow)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
