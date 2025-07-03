package service

import (
	"context"

	"github.com/dopp1e/webingo/backend/internal/dto"
	"github.com/dopp1e/webingo/backend/internal/errors"
	"github.com/dopp1e/webingo/backend/internal/model"
	"github.com/dopp1e/webingo/backend/internal/store"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type BoardService struct {
	db *gorm.DB
}

func (s *BoardService) CreateBoard(ctx context.Context, request *dto.BoardCreateRequest, userID uuid.UUID) (*model.Board, error) {
	var board *model.Board
	err := store.WithTransaction(ctx, s.db, func(tx store.TransactionalStorage) error {
		var preparedSpaces []model.Space
		for _, spaceRequest := range request.Spaces {
			spaceModel := &model.Space{Content: spaceRequest.Content}
			foundOrCreatedSpace, err := tx.Spaces().FindOrCreateByContent(ctx, spaceModel)
			if err != nil {
				return err
			}
			preparedSpaces = append(preparedSpaces, *foundOrCreatedSpace)
		}

		b := &model.Board{
			Name:        request.Name,
			Description: request.Description,
			Tags:        pq.StringArray(request.Tags),
			Version:     1,
			UserID:      userID,
			Spaces:      preparedSpaces,
		}

		if err := tx.Boards().Create(ctx, b); err != nil {
			return err
		}

		board = b

		return nil
	})

	if err != nil {
		return nil, err
	}

	return board, nil
}

func (s *BoardService) GetById(ctx context.Context, boardId uuid.UUID) (*model.Board, error) {
	var board *model.Board
	err := store.WithTransaction(ctx, s.db, func(tx store.TransactionalStorage) error {
		b, err := tx.Boards().GetById(ctx, boardId)
		if err != nil {
			return err
		}
		board = b
		return nil
	})

	if err != nil {
		return nil, err
	}

	if board == nil {
		return nil, errors.ErrNotFound // No board found
	}

	return board, nil
}

func (s *BoardService) UpdateBoard(ctx context.Context, request *dto.BoardUpdateRequest, boardID uuid.UUID) (*model.Board, error) {
	var board *model.Board
	err := store.WithTransaction(ctx, s.db, func(tx store.TransactionalStorage) error {
		existingBoard, err := tx.Boards().GetById(ctx, boardID)
		if err != nil {
			return err
		}
		if existingBoard == nil {
			return errors.ErrNotFound
		}

		if existingBoard.Version != request.Version {
			return errors.ErrDataVersionMismatch
		}

		var preparedSpaces []model.Space
		for _, spaceRequest := range request.Spaces {
			spaceModel := &model.Space{Content: spaceRequest.Content}
			foundOrCreatedSpace, err := tx.Spaces().FindOrCreateByContent(ctx, spaceModel)
			if err != nil {
				return err
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
			Tags:        pq.StringArray(request.Tags),
			Spaces:      preparedSpaces,
		}

		board.Spaces = preparedSpaces
		board.Version++ // increment the version for optimistic concurrency control

		if err := tx.Boards().Update(ctx, board); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return board, nil
}

func (s *BoardService) DeleteBoard(ctx context.Context, boardId uuid.UUID) error {
	return store.WithTransaction(ctx, s.db, func(tx store.TransactionalStorage) error {
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

		return nil
	})
}

func (s *BoardService) AddComment(ctx context.Context, boardId uuid.UUID, userId uuid.UUID, request *dto.CommentPutRequest) (*model.BoardComment, error) {
	var comment *model.BoardComment
	err := store.WithTransaction(ctx, s.db, func(tx store.TransactionalStorage) error {
		exists, err := tx.Boards().Exists(ctx, boardId)
		if err != nil {
			return err
		}
		if !exists {
			return errors.ErrNotFound
		}

		c := &model.BoardComment{
			Content: request.Content,
			BoardID: boardId,
			UserID:  userId,
		}

		if err := tx.BoardComments().Create(ctx, c); err != nil {
			return err
		}

		comment = c
		return nil
	})

	if err != nil {
		return nil, err
	}

	return comment, nil
}
