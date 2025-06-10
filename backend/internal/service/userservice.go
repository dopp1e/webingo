package service

import (
	"context"

	"github.com/dopp1e/webingo/backend/internal/model"
	"github.com/dopp1e/webingo/backend/internal/store"
	"gorm.io/gorm"
)

type UserService struct {
	storage store.Storage
	db      *gorm.DB
}

func (s *UserService) Create(context.Context, *model.User) error {
	err := s.storage.Users.Create(context.Background(), &model.User{})
	return err
}

func (s *UserService) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	user, err := s.storage.Users.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Exists(ctx context.Context, username string) (bool, error) {
	exists, err := s.storage.Users.Exists(ctx, username)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *UserService) CreateIfNotExists(ctx context.Context, user *model.User) (*model.User, error) {
	exists, err := s.storage.Users.Exists(ctx, user.Username)

	if err != nil {
		return nil, err
	}

	if exists {
		found, err := s.GetByUsername(ctx, user.Username)
		if err != nil {
			return nil, err
		}
		return found, nil
	}

	if err := s.storage.Users.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
