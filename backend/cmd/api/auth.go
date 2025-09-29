package main

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/dopp1e/webingo/backend/internal/dto"
	"github.com/dopp1e/webingo/backend/internal/errors"
	"github.com/google/uuid"
)

// registerUserHandler godoc
//
//	@Summary		Registers a new user
//	@Description	Registers a new user with the provided username, email, and password.
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			user	body		dto.UserCreateRequest	true	"User registration details"
//	@Success		201		{object}	dto.UserActivateData	"User successfully registered"
//	@Failure		400		{object}	error					"Bad request"
//	@Failure		500		{object}	error					"Internal server error"
//	@Router			/authentication/register [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload dto.UserCreateRequest
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.validator.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	plainToken := uuid.New().String()
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	if err := app.service.Users.CreateAndInvite(r.Context(), payload, hashToken, app.config.mail.exp); err != nil {
		switch err {
		case errors.ErrUsernameExists, errors.ErrEmailExists, errors.ErrAlreadyExists:
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	userData := &dto.UserActivateData{
		Username: payload.Username,
		Email:    payload.Email,
		Token:    plainToken,
	}

	if err := app.jsonResponse(w, http.StatusCreated, userData); err != nil {
		app.internalServerError(w, r, err)
	}
}
