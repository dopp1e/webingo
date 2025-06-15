package service

import (
	"context"

	"github.com/dopp1e/webingo/backend/internal/model"
	"github.com/dopp1e/webingo/backend/internal/store"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FeedService struct {
	db *gorm.DB
}

func (s *FeedService) GetFeed(ctx context.Context, userId uuid.UUID) ([]model.Board, error) {
	tx, err := store.StartTransaction(ctx, s.db)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // defer rollback in case of error

	followedUsers, err := tx.Users().GetFollowedUserIDs(ctx, userId)
	if err != nil {
		return nil, err
	}

	boards, err := tx.Boards().GetByCreatorIDs(ctx, followedUsers)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return boards, nil
}
