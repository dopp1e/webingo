package store

import (
	"context"

	"github.com/dopp1e/webingo/backend/internal/model"
	"gorm.io/gorm"
)

type RoleStore struct {
	db *gorm.DB
}

func (s *RoleStore) Create(ctx context.Context, role *model.Role) error {
	result := s.db.WithContext(ctx).Create(role)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *RoleStore) GetByName(ctx context.Context, name string) (*model.Role, error) {
	var role model.Role
	result := s.db.WithContext(ctx).Where("name = ?", name).First(&role)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // No role found
		}
		return nil, result.Error // Other error
	}

	return &role, nil
}

func (s *RoleStore) Exists(ctx context.Context, name string) (bool, error) {
	var count int64
	result := s.db.WithContext(ctx).Model(&model.Role{}).Where("name = ?", name).Count(&count)

	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}
