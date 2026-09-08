package app

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/m-krasilnikov/urlshortener/internal/handler"
)

func NewRouter(storage handler.URLStorage, baseURL string) http.Handler {
	log.Println("NewRouter started")

	h := handler.New(storage, baseURL)

	log.Println("Handler created")

	router := mux.NewRouter()

	log.Println("mux router created")

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	})

	router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	})

	router.HandleFunc(
		"/",
		h.CreateShortURL,
	).Methods(http.MethodPost)

	log.Println("POST route created")

	router.HandleFunc(
		"/{id}",
		h.GetOriginalURL,
	).Methods(http.MethodGet)

	log.Println("GET route created")

	return router
}
