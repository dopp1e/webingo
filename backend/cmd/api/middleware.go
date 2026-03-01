package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/dopp1e/webingo/backend/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.unauthorizedResponse(w, r, fmt.Errorf("authorization header is missing or empty"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.unauthorizedResponse(w, r, fmt.Errorf("authorization header is malformed"))
			return
		}

		tokenStr := parts[1]
		jwtToken, err := app.auth.ValidateToken(tokenStr)
		if err != nil {
			app.unauthorizedResponse(w, r, err)
			return
		}

		claims, _ := jwtToken.Claims.(jwt.MapClaims)
		userIDstring, ok := claims["sub"].(string)
		if !ok {
			app.unauthorizedResponse(w, r, fmt.Errorf("invalid token claims"))
			return
		}

		userID, err := uuid.Parse(userIDstring)
		if err != nil {
			app.unauthorizedResponse(w, r, fmt.Errorf("invalid user ID in token claims"))
			return
		}

		user, err := app.service.Users.GetByID(r.Context(), userID)
		if err != nil {
			switch err {
			case gorm.ErrRecordNotFound:
				app.unauthorizedResponse(w, r, fmt.Errorf("user not found"))
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		if !user.IsActive {
			app.unauthorizedResponse(w, r, fmt.Errorf("user is not active"))
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, model.Context.User, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unauthorizedBasicResponse(w, r, fmt.Errorf("authorization header is missing or empty"))
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				app.unauthorizedBasicResponse(w, r, fmt.Errorf("authorization header is malformed"))
				return
			}

			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unauthorizedBasicResponse(w, r, err)
				return
			}

			username := app.config.auth.basic.user
			password := app.config.auth.basic.pass

			creds := strings.SplitN(string(decoded), ":", 2)

			if len(creds) != 2 || creds[0] != username || creds[1] != password {
				app.unauthorizedBasicResponse(w, r, fmt.Errorf("wrong credentials"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
