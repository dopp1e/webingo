package store

import (
	"context"

	"github.com/dopp1e/webingo/backend/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ActivityStore struct {
	db *gorm.DB
}

func (s *ActivityStore) Create(ctx context.Context, activity *model.Activity) error {
	result := s.db.WithContext(ctx).Create(activity)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *ActivityStore) GetPaginatedActivitiesByActors(ctx context.Context, actorIDs []uuid.UUID, fq model.PaginatedFeedQuery) ([]model.Activity, error) {
	var activities []model.Activity
	query := s.db.WithContext(ctx).
		Where("actor_id IN ?", actorIDs).
		Preload("Actor").
		Offset(fq.Offset).
		Limit(fq.Limit)

	if fq.Since != "" {
		query = query.Where("updated_at > ?", fq.Since)
	}

	if fq.Until != "" {
		query = query.Where("updated_at < ?", fq.Until)
	}

	order := "created_at desc"
	if fq.Sort == "asc" {
		order = "created_at asc"
	} else if fq.Sort != "desc" {
		return nil, gorm.ErrInvalidField
	}
	query = query.Order(order)

	result := query.Find(&activities)
	if result.Error != nil {
		return nil, result.Error
	}

	return activities, nil
}
