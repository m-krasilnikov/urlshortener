package main

import (
	"log"
	"net/http"

	"github.com/m-krasilnikov/urlshortener/internal/app"
	"github.com/m-krasilnikov/urlshortener/internal/config"
)

func main() {
	cfg := config.New()

	router := app.NewRouter(cfg.BaseURL)

	log.Printf(
		"server started at http://%s",
		cfg.ServerAddress,
	)

	if err := http.ListenAndServe(
		cfg.ServerAddress,
		router,
	); err != nil {
		log.Fatal(err)
	}
}
