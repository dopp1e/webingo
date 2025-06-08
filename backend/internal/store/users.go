package store

import (
	"context"

	"github.com/dopp1e/webingo/backend/internal/model"
	"gorm.io/gorm"
)

type UsersStore struct {
	db *gorm.DB
}

func (s *UsersStore) Create(ctx context.Context, user *model.User) error {
	result := s.db.WithContext(ctx).Create(user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *UsersStore) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	result := s.db.WithContext(ctx).Where("username = ?", username).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil // User found
}

func (s *UsersStore) Exists(ctx context.Context, username string) (bool, error) {
	var count int64
	result := s.db.WithContext(ctx).Model(&model.User{}).Where("username = ?", username).Count(&count)

	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil // Return true if user exists
}
