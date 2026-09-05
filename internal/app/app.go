package app

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/m-krasilnikov/urlshortener/internal/handler"
	"github.com/m-krasilnikov/urlshortener/internal/storage"
)

func NewRouter(baseURL string) http.Handler {
	storage := storage.NewMemoryStorage()

	h := handler.New(
		storage,
		baseURL,
	)

	router := mux.NewRouter()

	router.HandleFunc(
		"/",
		h.CreateShortURL,
	).Methods(http.MethodPost)

	router.HandleFunc(
		"/{id}",
		h.GetOriginalURL,
	).Methods(http.MethodGet)

	router.NotFoundHandler = http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.Error(
				w,
				"Bad Request",
				http.StatusBadRequest,
			)
		},
	)

	router.MethodNotAllowedHandler = http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.Error(
				w,
				"Bad Request",
				http.StatusBadRequest,
			)
		},
	)

	return router
}
