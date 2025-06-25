package main

import (
	"net/http"

	"github.com/dopp1e/webingo/backend/internal/errors"
	"github.com/dopp1e/webingo/backend/internal/model"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	fq := model.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}

	fq, err := fq.Parse(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.validator.Struct(fq); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.getUserFromContext(r.Context())
	if err != nil {
		if err == errors.ErrNotFound {
			app.notFoundResponse(w, r, err)
			return
		} else {
			app.internalServerError(w, r, err)
			return
		}
	}

	feed, err := app.service.Feed.GetFeed(r.Context(), user.ID, fq)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
	}
}
