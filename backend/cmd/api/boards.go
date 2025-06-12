package main

import (
	"net/http"

	"github.com/dopp1e/webingo/backend/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

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
	boardID, err := uuid.Parse(chi.URLParam(r, "boardID"))

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	board, err := app.service.Boards.GetById(ctx, boardID)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			app.notFoundResponse(w, r, err)
		} else {
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := writeJSON(w, http.StatusOK, board); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) putBoardHandler(w http.ResponseWriter, r *http.Request) {
	boardID, err := uuid.Parse(chi.URLParam(r, "boardID"))

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// ensure that the board exists
	ctx := r.Context()
	if board, err := app.service.Boards.GetById(ctx, boardID); err != nil {
		if err == gorm.ErrRecordNotFound || board == nil {
			app.notFoundResponse(w, r, err)
			return
		}
		app.internalServerError(w, r, err)
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

	board := &model.Board{
		Base: model.Base{
			ID: boardID,
		},
		BoardPayload: payload,
	}

	if err := app.service.Boards.UpdateBoard(ctx, board); err != nil {
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
