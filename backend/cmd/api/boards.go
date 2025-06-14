package main

import (
	"context"
	"net/http"

	"github.com/dopp1e/webingo/backend/internal/dto"
	"github.com/dopp1e/webingo/backend/internal/errors"
	"github.com/dopp1e/webingo/backend/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (app *application) boardContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		boardId, err := uuid.Parse(chi.URLParam(r, "boardID"))
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		ctx := r.Context()
		board, err := app.service.Boards.GetById(ctx, boardId)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				app.notFoundResponse(w, r, err)
			} else {
				app.internalServerError(w, r, err)
			}
			return
		}
		if board == nil {
			app.notFoundResponse(w, r, errors.ErrNotFound)
			return
		}

		ctx = context.WithValue(ctx, model.Context.Board, board)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) createBoardHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: add verification of user making rqeust
	var payload dto.BoardCreateRequest
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.validator.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	board, err := app.service.Boards.CreateBoard(ctx, &payload, uuid.MustParse("558cf4e5-d324-4da7-a5eb-5836674ede97")) // Replace uuid.New() with the actual user ID from the context or session
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, board); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getBoardHandler(w http.ResponseWriter, r *http.Request) {
	board := getBoardFromContext(r)

	if board == nil {
		app.notFoundResponse(w, r, errors.ErrNotFound)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, board); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) putBoardHandler(w http.ResponseWriter, r *http.Request) {
	board := getBoardFromContext(r)
	if board == nil {
		app.notFoundResponse(w, r, errors.ErrNotFound)
		return
	}

	var payload dto.BoardUpdateRequest
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.validator.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	board, err := app.service.Boards.UpdateBoard(r.Context(), &payload, board.ID)
	if err != nil {
		switch err {
		case errors.ErrNotFound:
			app.notFoundResponse(w, r, err)
			return
		case errors.ErrDataVersionMismatch:
			app.conflictResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusOK, board); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) deleteBoardHandler(w http.ResponseWriter, r *http.Request) {
	boardID, err := uuid.Parse(chi.URLParam(r, "boardID"))

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	if err := app.service.Boards.DeleteBoard(ctx, boardID); err != nil {
		if err == gorm.ErrRecordNotFound {
			app.notFoundResponse(w, r, err)
		} else {
			app.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := app.jsonResponse(w, http.StatusOK, "board deleted successfully"); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) createBoardCommentHandler(w http.ResponseWriter, r *http.Request) {
	board := getBoardFromContext(r)
	if board == nil {
		app.notFoundResponse(w, r, errors.ErrNotFound)
		return
	}

	var payload dto.CommentPutRequest
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.validator.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	comment, err := app.service.Boards.AddComment(ctx, board.ID, uuid.MustParse("558cf4e5-d324-4da7-a5eb-5836674ede97"), &payload)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, comment); err != nil {
		app.internalServerError(w, r, err)
	}
}

func getBoardFromContext(r *http.Request) *model.Board {
	board, ok := r.Context().Value(model.Context.Board).(*model.Board)
	if !ok {
		return nil
	}
	return board
}
