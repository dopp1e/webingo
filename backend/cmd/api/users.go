package main

import (
	"context"
	"net/http"

	"github.com/dopp1e/webingo/backend/internal/errors"
	"github.com/dopp1e/webingo/backend/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// getUserHandler godoc
//
//	@Summary		Fetches a user profile
//	@Description	Retrieves the user profile based on the user ID
//	@Tags			users
//	@Accept			json
//
//	@Produce		json
//
//	@Param			userID	path		string	true	"User ID"	Format(uuid)
//	@Success		200		{object}	model.User
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//
//	@Security		ApiKeyAuth
//	@Router			/users/{userID} [get]
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

// FollowUser godoc
//
//	@Summary		Follow a user
//	@Description	Follows a user by their ID
//	@Tags			users
//	@Accept			json
//
//	@Produce		json
//	@Param			userID	path	string	true	"User ID"	Format(uuid)
//
//	@Success		204		"Followed successfully"
//
//	@Failure		400		{object}	error	"User payload missing or invalid"
//	@Failure		404		{object}	error	"User not found"
//	@Failure		409		{object}	error	"Already following this user"
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/follow [post]
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

// UnfollowUser godoc
//
//	@Summary		Unfollow a user
//	@Description	Unfollows a user by their ID
//	@Tags			users
//	@Accept			json
//
//	@Produce		json
//	@Param			userID	path	string	true	"User ID"	Format(uuid)
//	@Success		204		"Unfollowed successfully"
//	@Failure		400		{object}	error	"User payload missing or invalid"
//	@Failure		404		{object}	error	"User not found"
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/unfollow [post]
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

// ActivateUser godoc
//
//	@Summary		Activate a previously registered user
//	@Description	Activates a user account using the provided activation token
//	@Tags			users
//
//	@Produce		json
//	@Param			token	path		string	true	"Activation Token"
//	@Success		201		{string}	string	"User activated successfully"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/activate/{token} [put]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Implementation for user activation goes here
	token := chi.URLParam(r, "token")
	if token == "" {
		app.badRequestResponse(w, r, errors.ErrTokenRequired)
		return
	}

	if err := app.service.Users.Activate(r.Context(), token); err != nil {
		if err == errors.ErrNotFound {
			app.notFoundResponse(w, r, err)
			return
		} else {
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusCreated, "User activated successfully"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
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
