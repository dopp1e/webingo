package service

import (
	"context"

	"github.com/dopp1e/webingo/backend/internal/dto"
	"github.com/dopp1e/webingo/backend/internal/model"
	"github.com/dopp1e/webingo/backend/internal/store"
	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

func (s *AuthService) CreateUser(ctx context.Context, payload *dto.UserCreateRequest) error {
	tx, err := store.StartTransaction(ctx, s.db)
	if err != nil {
		return err
	}
	defer tx.Rollback() // defer rollback in case of error
	user := &model.User{
		Username: payload.Username,
		Email:    payload.Email,
	}
	user.SetPassword(payload.Password)
	return nil
}
