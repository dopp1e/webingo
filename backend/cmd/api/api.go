package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dopp1e/webingo/backend/docs"
	"github.com/dopp1e/webingo/backend/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

type application struct {
	config    config
	service   service.Service
	validator validator.Validate
	logger    *zap.SugaredLogger
}

type config struct {
	dsn    string
	db     dbConfig
	env    string
	apiUrl string
}

type dbConfig struct {
	dsn          string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.dsn)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		r.Route("/boards", func(r chi.Router) {
			r.Put("/", app.createBoardHandler)

			r.Route("/{boardID}", func(r chi.Router) {
				r.Use(app.boardContextMiddleware)

				r.Get("/", app.getBoardHandler)
				r.Put("/", app.putBoardHandler)
				r.Delete("/", app.deleteBoardHandler)

				r.Route("/comments", func(r chi.Router) {
					r.Put("/", app.createBoardCommentHandler)
				})
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.userContextMiddleware)

				r.Get("/", app.getUserHandler)

				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})

			r.Group(func(r chi.Router) {
				r.Get("/feed", app.getUserFeedHandler)
			})
		})

		// Public routes
		r.Route("/authentication", func(r chi.Router) {
			r.Post("/register", app.registerUserHandler)
		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	//
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.apiUrl
	docs.SwaggerInfo.BasePath = "/v1"

	srv := http.Server{
		Addr:         app.config.dsn,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	app.logger.Infow("Server has started", "addr", app.config.dsn, "env", app.config.env)

	return srv.ListenAndServe()
}
