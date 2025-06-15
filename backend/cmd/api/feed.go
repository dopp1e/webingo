package main

import (
	"net/http"

	"github.com/dopp1e/webingo/backend/internal/errors"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
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

	feed, err := app.service.Feed.GetFeed(r.Context(), user.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
	}
}
