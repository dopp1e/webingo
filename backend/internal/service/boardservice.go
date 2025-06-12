package service

import (
	"context"

	"github.com/dopp1e/webingo/backend/internal/model"
	"github.com/dopp1e/webingo/backend/internal/store"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BoardService struct {
	db *gorm.DB
}

func (s *BoardService) CreateBoard(ctx context.Context, board *model.Board) error {
	tx, err := store.StartTransaction(ctx, s.db)
	if err != nil {
		return err
	}
	defer tx.Rollback() // defer rollback in case of error

	var preparedSpaces []model.Space
	for _, incomingSpace := range board.Spaces {
		foundOrCreatedSpace, err := tx.Spaces().FindOrCreateByContent(ctx, &incomingSpace)
		if err != nil {
			return err
		}
		preparedSpaces = append(preparedSpaces, *foundOrCreatedSpace)
	}

	board.Spaces = preparedSpaces

	if err := tx.Boards().Create(ctx, board); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *BoardService) GetById(ctx context.Context, boardId uuid.UUID) (*model.Board, error) {
	tx, err := store.StartTransaction(ctx, s.db)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // defer rollback in case of error

	board, err := tx.Boards().GetById(ctx, boardId)
	if err != nil {
		return nil, err
	}

	return board, nil
}

func (s *BoardService) UpdateBoard(ctx context.Context, board *model.Board) error {
	tx, err := store.StartTransaction(ctx, s.db)
	if err != nil {
		return err
	}
	defer tx.Rollback() // defer rollback in case of error

	var preparedSpaces []model.Space
	for _, incomingSpace := range board.Spaces {
		foundOrCreatedSpace, err := tx.Spaces().FindOrCreateByContent(ctx, &incomingSpace)
		if err != nil {
			return err
		}
		preparedSpaces = append(preparedSpaces, *foundOrCreatedSpace)
	}

	board.Spaces = preparedSpaces

	if err := tx.Boards().Update(ctx, board); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *BoardService) DeleteBoard(ctx context.Context, boardId uuid.UUID) error {
	tx, err := store.StartTransaction(ctx, s.db)
	if err != nil {
		return err
	}
	defer tx.Rollback() // defer rollback in case of error

	board, err := tx.Boards().GetById(ctx, boardId)
	if err != nil {
		return err
	}

	if board == nil {
		return gorm.ErrRecordNotFound
	}

	if err := tx.Model(board).Association("Spaces").Clear(); err != nil {
		return err
	}

	err = tx.Boards().Delete(ctx, board)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
