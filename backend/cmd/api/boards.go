package main

import (
	"net/http"

	"github.com/dopp1e/webingo/backend/internal/model"
)

func (app *application) createBoardHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: add verification of user making rqeust
	var payload model.BoardPayload
	if err := readJSON(w, r, &payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	board := &model.Board{
		BoardPayload: payload,
	}

	ctx := r.Context()

	if err := app.store.Boards.Create(ctx, board); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeJSON(w, http.StatusCreated, board); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
	}
}
