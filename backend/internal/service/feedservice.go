package service

import (
	"context"

	"github.com/dopp1e/webingo/backend/internal/dto"
	"github.com/dopp1e/webingo/backend/internal/model"
	"github.com/dopp1e/webingo/backend/internal/store"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FeedService struct {
	db *gorm.DB
}

func (s *FeedService) GetFeed(ctx context.Context, userId uuid.UUID, fq model.PaginatedFeedQuery) ([]dto.FeedItem, error) {
	tx, err := store.StartTransaction(ctx, s.db)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // defer rollback in case of error

	followedUsers, err := tx.Users().GetFollowedUserIDs(ctx, userId)
	if err != nil {
		return nil, err
	}

	activities, err := tx.Activities().GetPaginatedActivitiesByActors(ctx, followedUsers, fq.Offset, fq.Limit, fq.Sort)
	if err != nil {
		return nil, err
	}

	feedItems := make([]dto.FeedItem, 0, len(activities))
	for _, activity := range activities {
		feedItem := dto.FeedItem{
			ID:          activity.ID,
			Actor:       dto.FeedUserData{ID: activity.Actor.ID, Username: activity.Actor.Username},
			Verb:        activity.Verb,
			CreatedAt:   activity.CreatedAt,
			TargetID:    activity.TargetID,
			TargetType:  activity.TargetType,
			ContextID:   activity.ContextID,
			ContextType: activity.ContextType,
		}

		switch activity.TargetType {
		case model.ActivityTargetTypeBoard:
			board, boardErr := tx.Boards().GetById(ctx, activity.TargetID)
			if boardErr != nil {
				if boardErr == gorm.ErrRecordNotFound {
					continue // Skip if board not found
				}
				return nil, boardErr
			}
			feedItem.Target = dto.BoardFeedItem{
				ID:          board.ID,
				Name:        board.Name,
				Description: board.Description,
				Tags:        board.Tags,
			}
		case model.ActivityTargetTypeComment:
			comment, commentErr := tx.BoardComments().GetByID(ctx, activity.TargetID)
			if commentErr != nil {
				continue
			}
			feedItem.Target = dto.CommentFeedItem{
				ID:      comment.ID,
				Content: comment.Content,
			}
		case model.ActivityTargetTypeGame:
			game, gameErr := tx.Games().GetByID(ctx, activity.TargetID)
			if gameErr != nil {
				continue
			}
			markedFieldsCount, markedFieldsErr := tx.Games().GetGameMarkedFieldCount(ctx, game.ID)
			if markedFieldsErr != nil {
				continue
			}
			feedItem.Target = dto.GameFeedItem{
				ID:           game.ID,
				YoutubeID:    game.YoutubeVideoID,
				Height:       game.Height,
				Width:        game.Width,
				MarkedFields: markedFieldsCount,
			}
		}

		feedItems = append(feedItems, feedItem)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return feedItems, nil
}
