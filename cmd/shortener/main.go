package main

import (
	"log"
	"net/http"

	"github.com/m-krasilnikov/urlshortener/internal/app"
	"github.com/m-krasilnikov/urlshortener/internal/config"
	"github.com/m-krasilnikov/urlshortener/internal/storage"
)

func main() {
	cfg := config.New()
	st := storage.NewMemoryStorage()
	router := app.NewRouter(st, cfg.BaseURL)
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
