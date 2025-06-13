package service

import (
	"context"

	"github.com/dopp1e/webingo/backend/internal/dto"
	"github.com/dopp1e/webingo/backend/internal/errors"
	"github.com/dopp1e/webingo/backend/internal/model"
	"github.com/dopp1e/webingo/backend/internal/store"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BoardService struct {
	db *gorm.DB
}

func (s *BoardService) CreateBoard(ctx context.Context, request *dto.BoardCreateRequest, userID uuid.UUID) (*model.Board, error) {
	tx, err := store.StartTransaction(ctx, s.db)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // defer rollback in case of error

	var preparedSpaces []model.Space
	for _, spaceRequest := range request.Spaces {
		spaceModel := &model.Space{Content: spaceRequest.Content}
		foundOrCreatedSpace, err := tx.Spaces().FindOrCreateByContent(ctx, spaceModel)
		if err != nil {
			return nil, err
		}
		preparedSpaces = append(preparedSpaces, *foundOrCreatedSpace)
	}

	board := &model.Board{
		Name:        request.Name,
		Description: request.Description,
		Tags:        request.Tags,
		Version:     1,
		UserID:      userID,
		Spaces:      preparedSpaces,
	}

	if err := tx.Boards().Create(ctx, board); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return board, nil
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

	if board == nil {
		return nil, errors.ErrNotFound // No board found
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return board, nil
}

func (s *BoardService) UpdateBoard(ctx context.Context, request *dto.BoardUpdateRequest, boardID uuid.UUID) (*model.Board, error) {
	tx, err := store.StartTransaction(ctx, s.db)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // defer rollback in case of error

	// ensure the board version matches the one in the database
	existingBoard, err := tx.Boards().GetById(ctx, boardID)
	if err != nil {
		return nil, err
	}
	if existingBoard == nil {
		return nil, errors.ErrNotFound
	}

	if existingBoard.Version != request.Version {
		return nil, errors.ErrDataVersionMismatch
	}

	var preparedSpaces []model.Space
	for _, spaceRequest := range request.Spaces {
		spaceModel := &model.Space{Content: spaceRequest.Content}
		foundOrCreatedSpace, err := tx.Spaces().FindOrCreateByContent(ctx, spaceModel)
		if err != nil {
			return nil, err
		}
		preparedSpaces = append(preparedSpaces, *foundOrCreatedSpace)
	}

	board := &model.Board{
		Base: model.Base{
			ID: boardID,
		},
		Name:        request.Name,
		Description: request.Description,
		Version:     request.Version + 1,
		Tags:        request.Tags,
		Spaces:      preparedSpaces,
	}

	board.Spaces = preparedSpaces
	board.Version++ // increment the version for optimistic concurrency control

	if err := tx.Boards().Update(ctx, board); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return board, nil
}

func (s *BoardService) DeleteBoard(ctx context.Context, boardId uuid.UUID) error {
	tx, err := store.StartTransaction(ctx, s.db)
	if err != nil {
		return err
	}
	defer tx.Rollback() // defer rollback in case of error

	request, err := tx.Boards().GetById(ctx, boardId)
	if err != nil {
		return err
	}

	if request == nil {
		return errors.ErrNotFound
	}

	if err := tx.Model(request).Association("Spaces").Clear(); err != nil {
		return err
	}

	err = tx.Boards().Delete(ctx, request)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
