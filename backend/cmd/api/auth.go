package main

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/dopp1e/webingo/backend/internal/dto"
	"github.com/dopp1e/webingo/backend/internal/errors"
	"github.com/dopp1e/webingo/backend/internal/mailer"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
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

	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      payload.Username,
		ActivationURL: app.config.frontendUrl + "/activate?token=" + plainToken,
	}

	err := app.mailer.Send(
		mailer.UserInvitationTemplate,
		payload.Username,
		payload.Email,
		vars,
		app.config.env != "production",
	)

	if err != nil {
		app.logger.Errorw("failed to send invitation email", "error", err)

		// rollback user creation (SAGA pattern)
		if delErr := app.service.Users.DeleteByUsername(r.Context(), payload.Username); delErr != nil {
			app.logger.Errorw("failed to rollback user creation after email failure", "error", delErr)
		}

		app.internalServerError(w, r, err)
		return
	}

	app.logger.Info("sent invitation email to " + payload.Email)

	if err := app.jsonResponse(w, http.StatusCreated, userData); err != nil {
		app.internalServerError(w, r, err)
	}
}

// createTokenHandler godoc
//
//	@Summary		Creates an authentication token
//	@Description	Creates a JWT authentication token for a user with the provided username and password.
//	@Tags			authentication
//
//	@Accept			json
//	@Produce		json
//	@Param			payload	body	dto.CreateUserTokenRequest	true	"User Credentials"
//	@Success		200		body	{string}					"Token"
//	@Failure		400		body	{object}					error
//	@Failure		401		body	{object}					error
//	@Failure		500		body	{object}					error
//	@Router			/authentication/token [post]
func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {
	var payload dto.CreateUserTokenRequest
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.validator.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.service.Users.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound || err == errors.ErrNotFound {
			app.unauthorizedResponse(w, r, err)
			return
		}
		app.internalServerError(w, r, err)
		return
	}

	ok := user.CheckPassword(payload.Password)
	if !ok {
		app.unauthorizedResponse(w, r, errors.ErrIncorrectPassword)
		return
	}

	claims := jwt.MapClaims{
		"sub": user.ID.String(),
		"aud": app.config.auth.token.issuer,
		"iss": app.config.auth.token.issuer,
		"exp": time.Now().Add(app.config.auth.token.expiration).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
	}

	token, err := app.auth.GenerateToken(claims)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, token); err != nil {
		app.internalServerError(w, r, err)
	}
}
