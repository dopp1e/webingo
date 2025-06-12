package main

import (
	"context"
	"net/http"

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
			app.notFoundResponse(w, r, gorm.ErrRecordNotFound)
			return
		}

		ctx = context.WithValue(ctx, model.BoardCtx, board)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) createBoardHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: add verification of user making rqeust
	var payload model.BoardPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.validator.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	board := &model.Board{
		BoardPayload: payload,
	}

	ctx := r.Context()

	if err := app.service.Boards.CreateBoard(ctx, board); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusCreated, board); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getBoardHandler(w http.ResponseWriter, r *http.Request) {
	board := getBoardFromContext(r)

	if board == nil {
		app.notFoundResponse(w, r, gorm.ErrRecordNotFound)
		return
	}

	if err := writeJSON(w, http.StatusOK, board); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) putBoardHandler(w http.ResponseWriter, r *http.Request) {
	board := getBoardFromContext(r)
	if board == nil {
		app.notFoundResponse(w, r, gorm.ErrRecordNotFound)
		return
	}

	var payload model.BoardPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.validator.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	board = &model.Board{
		Base: model.Base{
			ID: board.ID,
		},
		BoardPayload: payload,
	}

	if err := app.service.Boards.UpdateBoard(r.Context(), board); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusOK, board); err != nil {
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
	if err := writeJSON(w, http.StatusOK, map[string]string{"message": "Board deleted successfully"}); err != nil {
		app.internalServerError(w, r, err)
	}
}

func getBoardFromContext(r *http.Request) *model.Board {
	board, ok := r.Context().Value(model.BoardCtx).(*model.Board)
	if !ok {
		return nil
	}
	return board
}
