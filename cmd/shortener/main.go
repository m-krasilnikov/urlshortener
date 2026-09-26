package main

import (
	"log"
	"net/http"

	"go.uber.org/zap"

	"github.com/m-krasilnikov/urlshortener/internal/app"
	"github.com/m-krasilnikov/urlshortener/internal/config"
	"github.com/m-krasilnikov/urlshortener/internal/middleware"
	"github.com/m-krasilnikov/urlshortener/internal/storage"
)

func main() {
	cfg := config.New()
	st := storage.NewMemoryStorage()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	router := app.NewRouter(st, cfg.BaseURL)

	handler := middleware.WithLogging(logger)(router)

	log.Printf(
		"server started at http://%s",
		cfg.ServerAddress,
	)

	if err := http.ListenAndServe(
		cfg.ServerAddress,
		handler,
	); err != nil {
		log.Fatal(err)
	}
}
