package service

import (
	"context"

	"github.com/dopp1e/webingo/backend/internal/errors"
	"github.com/dopp1e/webingo/backend/internal/model"
	"github.com/dopp1e/webingo/backend/internal/store"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func (s *UserService) Create(ctx context.Context, user *model.User) error {
	tx, err := store.StartTransaction(ctx, s.db)
	if err != nil {
		return err
	}
	defer tx.Rollback() // defer rollback in case of error

	err = tx.Users().Create(context.Background(), user)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return err
}

func (s *UserService) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	tx, err := store.StartTransaction(ctx, s.db)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // defer rollback in case of error
	user, err := tx.Users().GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.ErrNotFound // No user found
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Exists(ctx context.Context, username string) (bool, error) {
	tx, err := store.StartTransaction(ctx, s.db)
	if err != nil {
		return false, err
	}
	defer tx.Rollback() // defer rollback in case of error

	exists, err := tx.Users().Exists(ctx, username)
	if err != nil {
		return false, err
	}

	if err := tx.Commit(); err != nil {
		return false, err
	}

	return exists, nil
}

func (s *UserService) CreateIfNotExists(ctx context.Context, user *model.User) (*model.User, error) {
	tx, err := store.StartTransaction(ctx, s.db)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // defer rollback in case of error

	exists, err := tx.Users().Exists(ctx, user.Username)
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

	if err := tx.Users().Create(ctx, user); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return user, nil
}
