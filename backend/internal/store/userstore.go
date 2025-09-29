package store

import (
	"context"

	"github.com/dopp1e/webingo/backend/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserStore struct {
	db *gorm.DB
}

func (s *UserStore) Create(ctx context.Context, user *model.User) error {
	result := s.db.WithContext(ctx).Create(user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *UserStore) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	result := s.db.WithContext(ctx).Where("username = ?", username).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil // User found
}

func (s *UserStore) Exists(ctx context.Context, username string) (bool, error) {
	var count int64
	result := s.db.WithContext(ctx).Model(&model.User{}).Where("username = ?", username).Count(&count)

	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil // Return true if user exists
}

func (s *UserStore) GetByID(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	var user model.User
	result := s.db.WithContext(ctx).Where("id = ?", userID).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil // User found
}

func (s *UserStore) GetFollowedUserIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	var followedUsers []uuid.UUID
	result := s.db.WithContext(ctx).
		Model(&model.Follow{}).
		Where("follower_id = ?", userID).
		Pluck("followed_id", &followedUsers)

	if result.Error != nil {
		return nil, result.Error
	}

	return followedUsers, nil // Return list of followed user IDs
}

func (s *UserStore) Activate(ctx context.Context, userID uuid.UUID) error {
	result := s.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Update("is_active", true)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
