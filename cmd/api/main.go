package main

import (
	"log"

	"github.com/dopp1e/webingo/backend/internal/env"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
	}

	app := &application{
		cfg,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}