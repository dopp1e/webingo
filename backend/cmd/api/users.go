package main

import (
	"context"
	"net/http"

	"github.com/dopp1e/webingo/backend/internal/errors"
	"github.com/dopp1e/webingo/backend/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
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

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followedUser, err := app.getUserFromContext(r.Context())
	if err != nil {
		if err == errors.ErrNotFound {
			app.notFoundResponse(w, r, err)
			return
		} else {
			app.internalServerError(w, r, err)
			return
		}
	}

	followerID := uuid.MustParse("558cf4e5-d324-4da7-a5eb-5836674ede97") // Replace with actual user ID from context or session

	err = app.service.Users.FollowUser(r.Context(), followerID, followedUser.ID)
	if err != nil {
		if err == errors.ErrAlreadyExists {
			app.conflictResponse(w, r, err)
			return
		} else {
			app.internalServerError(w, r, err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	followedUser, err := app.getUserFromContext(r.Context())
	if err != nil {
		if err == errors.ErrNotFound {
			app.notFoundResponse(w, r, err)
			return
		} else {
			app.internalServerError(w, r, err)
			return
		}
	}

	followerID := uuid.MustParse("558cf4e5-d324-4da7-a5eb-5836674ede97") // Replace with actual user ID from context or session

	err = app.service.Users.UnfollowUser(r.Context(), followerID, followedUser.ID)
	if err != nil {
		if err == errors.ErrNotFound {
			app.notFoundResponse(w, r, err)
			return
		} else {
			app.internalServerError(w, r, err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := uuid.Parse(chi.URLParam(r, "userID"))
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		user, err := app.service.Users.GetByID(r.Context(), userID)
		if err != nil {
			if err == errors.ErrNotFound {
				app.notFoundResponse(w, r, err)
				return
			} else {
				app.internalServerError(w, r, err)
				return
			}
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, model.Context.User, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) getUserFromContext(ctx context.Context) (*model.User, error) {
	user, ok := ctx.Value(model.Context.User).(*model.User)
	if !ok || user == nil {
		return nil, errors.ErrNotFound
	}
	return user, nil
}
