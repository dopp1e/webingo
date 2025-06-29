package main

import (
	"net/http"

	"github.com/dopp1e/webingo/backend/internal/errors"
	"github.com/dopp1e/webingo/backend/internal/model"
)

// getUserFeedHandler godoc
//
//	@Summary		Get user feed
//	@Description	Retrieves the feed for the authenticated user
//	@Tags			feed
//	@Accept			json
//
//	@Produce		json
//	@Param			limit	query		int		false	"Number of items per page"
//	@Param			offset	query		int		false	"Offset for pagination"
//	@Param			sort	query		string	false	"Sort order (asc/desc)"
//	@Param			tags	query		string	false	"Comma-separated list of tags"
//	@Param			since	query		string	false	"Filter items created since this time" Format(date-time)
//	@Param			until	query		string	false	"Filter items created until this time" Format(date-time)
//
//	@Success		200		{object}	[]dto.FeedItem
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//
//	@Security		ApiKeyAuth
//	@Router			/feed [get]
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
