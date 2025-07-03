package service

import (
	"context"
	"time"

	"github.com/dopp1e/webingo/backend/internal/dto"
	"github.com/dopp1e/webingo/backend/internal/errors"
	"github.com/dopp1e/webingo/backend/internal/model"
	"github.com/dopp1e/webingo/backend/internal/store"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func (s *UserService) Create(ctx context.Context, req dto.UserCreateRequest) (*model.User, error) {
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
	}
	user.SetPassword(req.Password)
	err := store.WithTransaction(ctx, s.db, func(tx store.TransactionalStorage) error {
		err := tx.Users().Create(context.Background(), user)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user *model.User
	err := store.WithTransaction(ctx, s.db, func(tx store.TransactionalStorage) error {
		u, err := tx.Users().GetByUsername(ctx, username)
		if err != nil {
			return err
		}
		user = u
		return nil
	})

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.ErrNotFound // No user found
	}

	return user, nil
}

func (s *UserService) Exists(ctx context.Context, username string) (bool, error) {
	var exists bool
	err := store.WithTransaction(ctx, s.db, func(tx store.TransactionalStorage) error {
		e, err := tx.Users().Exists(ctx, username)

		if err != nil {
			return err
		}
		exists = e
		return nil
	})

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (s *UserService) CreateIfNotExists(ctx context.Context, user *model.User) (*model.User, error) {
	var foundUser *model.User
	err := store.WithTransaction(ctx, s.db, func(tx store.TransactionalStorage) error {
		exists, err := tx.Users().Exists(ctx, user.Username)
		if err != nil {
			return err
		}

		if exists {
			found, err := s.GetByUsername(ctx, user.Username)
			if err != nil {
				return err
			}
			foundUser = found
			return nil
		}

		if err := tx.Users().Create(ctx, user); err != nil {
			return err
		}

		foundUser = user

		return nil
	})

	if err != nil {
		return nil, err
	}

	return foundUser, nil
}

func (s *UserService) GetByID(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	var user *model.User
	err := store.WithTransaction(ctx, s.db, func(tx store.TransactionalStorage) error {
		u, err := tx.Users().GetByID(ctx, userID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.ErrNotFound // No user found
			}
			return err // Other error
		}
		user = u
		return nil
	})
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.ErrNotFound // No user found
	}
	return user, nil
}

func (s *UserService) FollowUser(ctx context.Context, followerID, followedUserID uuid.UUID) error {
	return store.WithTransaction(ctx, s.db, func(tx store.TransactionalStorage) error {
		follow := &model.Follow{
			FollowerID:     followerID,
			FollowedUserID: followedUserID,
		}

		err := tx.Follows().Create(ctx, follow)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.ErrNotFound // User not found
			} else if err == gorm.ErrDuplicatedKey {
				return errors.ErrAlreadyExists // Follow relationship already exists
			}
			return err // Other error
		}

		return nil
	})
}

func (s *UserService) UnfollowUser(ctx context.Context, followerID, followedUserID uuid.UUID) error {
	return store.WithTransaction(ctx, s.db, func(tx store.TransactionalStorage) error {
		follow := &model.Follow{
			FollowerID:     followerID,
			FollowedUserID: followedUserID,
		}

		err := tx.Follows().Delete(ctx, follow)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.ErrNotFound // Follow relationship not found
			}
			return err // Other error
		}

		return nil
	})
}

func (s *UserService) CreateAndInvite(ctx context.Context, req dto.UserCreateRequest, exp time.Duration) error {
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
	}
	user.SetPassword(req.Password)

	return store.WithTransaction(ctx, s.db, func(tx store.TransactionalStorage) error {
		if err := tx.Users().Create(ctx, user); err != nil {
			return err
		}

		invitation := &model.Invitation{
			UserID: user.ID,
			Token:  []byte(uuid.New().String()),
			Expiry: time.Now().Add(exp),
		}

		if err := tx.Invitations().Create(ctx, invitation); err != nil {
			return err
		}

		return nil
	})
}
