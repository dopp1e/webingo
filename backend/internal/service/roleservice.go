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
	tx, err := store.StartTransaction(ctx, s.db)
	if err != nil {
		return err
	}
	defer tx.Rollback() // defer rollback in case of error

	if err := tx.Roles().Create(ctx, role); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *RoleService) GetRoleByName(ctx context.Context, name string) (*model.Role, error) {
	tx, err := store.StartTransaction(ctx, s.db)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // defer rollback in case of error

	role, err := tx.Roles().GetByName(ctx, name)
	if err != nil {
		return nil, err
	}

	if role == nil {
		return nil, gorm.ErrRecordNotFound // No role found
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return role, nil
}

func (s *RoleService) CreateRoleIfNotExists(ctx context.Context, role *model.Role) (*model.Role, error) {
	tx, err := store.StartTransaction(ctx, s.db)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // defer rollback in case of error

	exists, err := tx.Roles().Exists(ctx, role.Name)
	if err != nil {
		return nil, err
	}

	if exists {
		// Role already exists, return it
		existingRole, err := tx.Roles().GetByName(ctx, role.Name)
		if err != nil {
			return nil, err
		}
		if existingRole == nil {
			return nil, errors.ErrNotFound // Should not happen if Exists is true
		}
		if err := tx.Commit(); err != nil {
			return nil, err
		}
		return existingRole, nil
	}

	if err := tx.Roles().Create(ctx, role); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return role, nil // New role created
}
