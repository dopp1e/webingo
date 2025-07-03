package service

import (
	"context"

	"github.com/dopp1e/webingo/backend/internal/errors"
	"github.com/dopp1e/webingo/backend/internal/model"
	"github.com/dopp1e/webingo/backend/internal/store"
	"gorm.io/gorm"
)

type RoleService struct {
	db *gorm.DB
}

func (s *RoleService) CreateRole(ctx context.Context, role *model.Role) error {
	return store.WithTransaction(ctx, s.db, func(tx store.TransactionalStorage) error {
		err := tx.Roles().Create(ctx, role)
		if err != nil {
			return err
		}
		return nil
	})
}

func (s *RoleService) GetRoleByName(ctx context.Context, name string) (*model.Role, error) {
	var role *model.Role
	err := store.WithTransaction(ctx, s.db, func(tx store.TransactionalStorage) error {
		r, err := tx.Roles().GetByName(ctx, name)
		if err != nil {
			return err
		}
		role = r
		return nil
	})

	if err != nil {
		return nil, err
	}

	if role == nil {
		return nil, errors.ErrNotFound // No role found
	}

	return role, nil
}

func (s *RoleService) CreateRoleIfNotExists(ctx context.Context, role *model.Role) (*model.Role, error) {
	var foundRole *model.Role
	err := store.WithTransaction(ctx, s.db, func(tx store.TransactionalStorage) error {
		exists, err := tx.Roles().Exists(ctx, role.Name)
		if err != nil {
			return err
		}

		if exists {
			// Role already exists, return it
			existingRole, err := tx.Roles().GetByName(ctx, role.Name)
			if err != nil {
				return err
			}
			foundRole = existingRole
			return nil
		}
		if err := tx.Roles().Create(ctx, role); err != nil {
			return err
		}
		foundRole = role
		return nil
	})

	if err != nil {
		return nil, err
	}
	if foundRole == nil {
		return nil, errors.ErrNotFound // No role found
	}
	return foundRole, nil
}
