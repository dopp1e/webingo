package main

import (
	"context"
	"time"

	"github.com/dopp1e/webingo/backend/internal/auth"
	"github.com/dopp1e/webingo/backend/internal/dto"
	"github.com/dopp1e/webingo/backend/internal/env"
	"github.com/dopp1e/webingo/backend/internal/mailer"
	"github.com/dopp1e/webingo/backend/internal/model"
	"github.com/dopp1e/webingo/backend/internal/service"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const version = "0.0.1"

//	@title			Webingo API
//	@description	API for Webingo, a bingo game application.

//	@license.name	GPL-3.0
//	@license.url	https://www.gnu.org/licenses/gpl-3.0.en.html

//	@BasePath	/v1

// @securityDefinitions.apiKey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	db_host := env.GetString("DB_HOST", "localhost")
	db_user := env.GetString("DB_USER", "admin")
	db_password := env.GetString("DB_PASSWORD", "adminpassword")
	db_name := env.GetString("DB_NAME", "webingo")
	db_port := env.GetString("DB_PORT", "5432")
	db_sslmode := env.GetString("DB_SSLMODE", "disable")
	db_timezone := env.GetString("DB_TIMEZONE", "UTC")
	dsn := env.GetString("DB_DSN", "host="+db_host+" user="+db_user+" password="+db_password+" dbname="+db_name+" port="+db_port+" sslmode="+db_sslmode+" TimeZone="+db_timezone)
	admin_username := env.GetString("ADMIN_USERNAME", "admin")
	admin_password := env.GetString("ADMIN_PASSWORD", "adminpassword")
	cfg := config{
		dsn:    env.GetString("ADDR", ":8080"),
		apiUrl: env.GetString("EXTERNAL_URL", "localhost:8080"),
		frontendUrl: env.GetString("FRONTEND_URL",
			"http://localhost:3000"),
		db: dbConfig{
			dsn:          dsn,
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "dev"),
		mail: mailConfig{
			exp: time.Hour * 24 * 3, // 3 days
			smtp: smtpConfig{
				host:     env.GetString("SMTP_HOST", "smtp.example.com"),
				port:     env.GetInt("SMTP_PORT", 587),
				username: env.GetString("SMTP_USERNAME", "user@example.com"),
				password: env.GetString("SMTP_PASSWORD", "password"),
				email:    env.GetString("SMTP_EMAIL", "noreply@example.com"),
			},
		},
		auth: authConfig{
			basic: basicAuthConfig{
				user: env.GetString("BASIC_AUTH_USER", "admin"),
				pass: env.GetString("BASIC_AUTH_PASSWORD", "password"),
			},
			token: tokenConfig{
				secret:          env.GetString("TOKEN_SECRET", "yoursuperdupersecretthingkeydoodad"),
				expiration:      time.Hour * time.Duration(env.GetInt("TOKEN_EXPIRATION", 24)),       // 1 day
				refreshInterval: time.Hour * time.Duration(env.GetInt("TOKEN_REFRESH_INTERVAL", 12)), // 12 hours
				issuer:          env.GetString("TOKEN_ISSUER", "webingo"),
			},
		},
	}

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	if cfg.env == "dev" {
		logger = zap.Must(zap.NewDevelopment()).Sugar()
	}

	defer logger.Sync() // flushes buffer, if any

	mailer := mailer.NewSMTPMailer(
		cfg.mail.smtp.host,
		cfg.mail.smtp.port,
		cfg.mail.smtp.username,
		cfg.mail.smtp.password,
		cfg.mail.smtp.email,
	)

	jwtAuthenticator := auth.NewJWTAuthenticator(cfg.auth.token.secret, cfg.auth.token.issuer, cfg.auth.token.issuer)

	// Database connection
	db, err := gorm.Open(postgres.Open(cfg.db.dsn), &gorm.Config{
		//DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		logger.Fatal(err)
	}

	if err := model.Migrate(db); err != nil {
		logger.Fatal(err)
	}

	logger.Info("Successfully connected to the database.")

	service := service.NewService(db)

	adminRoleRequest := dto.RoleCreateRequest{
		Name:        "admin",
		Description: "Administrator role with full access",
		Level:       1,
	}
	adminRole, err := service.Roles.CreateRoleIfNotExists(context.Background(), adminRoleRequest)
	if err != nil {
		logger.Fatal("failed to create admin role: %v", err)
	}
	userRoleRequest := dto.RoleCreateRequest{
		Name:        "user",
		Description: "Default user role with limited access",
		Level:       10,
	}
	_, err = service.Roles.CreateRoleIfNotExists(context.Background(), userRoleRequest)
	if err != nil {
		logger.Fatal("failed to create user role: %v", err)
	}
	adminUser := &model.User{
		Username: admin_username,
		Email:    "admin@example.com",
		IsActive: true,
		RoleID:   adminRole.ID,
	}
	adminUser.SetPassword(admin_password)
	adminUser, err = service.Users.CreateIfNotExists(context.Background(), adminUser)
	if err != nil {
		logger.Fatal("failed to create admin user: %v", err)
	}
	logger.Info("admin user available: %s with role %s", adminUser.Username, adminRole.Name)

	validate := validator.New(validator.WithRequiredStructEnabled())

	app := &application{
		config:    cfg,
		service:   *service,
		validator: *validate,
		logger:    logger,
		mailer:    mailer,
		auth:      jwtAuthenticator,
	}

	mux := app.mount()
	logger.Fatal(app.run(mux))
}
