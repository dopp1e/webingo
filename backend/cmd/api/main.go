package main

import (
	"context"
	"log"

	"github.com/dopp1e/webingo/backend/internal/env"
	"github.com/dopp1e/webingo/backend/internal/model"
	"github.com/dopp1e/webingo/backend/internal/store"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const version = "0.0.1"

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
		dsn: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			dsn:          dsn,
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "dev"),
	}

	db, err := gorm.Open(postgres.Open(cfg.db.dsn), &gorm.Config{
		//DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Panic(err)
	}

	if err := model.Migrate(db); err != nil {
		log.Panic(err)
	}

	log.Println("db connect")

	store := store.NewStorage(db)

	adminUser := &model.User{
		Username: admin_username,
		Password: admin_password,
	}

	adminExists, err := store.Users.Exists(context.Background(), adminUser.Username)
	if err != nil {
		log.Panicf("Failed to check if admin user exists: %v", err)
	}

	if adminExists {
		log.Println("Admin user already exists, skipping creation.")
	} else {
		if err := store.Users.Create(context.Background(), adminUser); err != nil {
			log.Panicf("Failed to create admin user: %v", err)
		}
		log.Println("Admin user created successfully.")
	}

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
