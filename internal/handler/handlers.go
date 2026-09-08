package handler

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/m-krasilnikov/urlshortener/internal/storage"
)

const idLength = 8

type URLStorage interface {
	Save(id, url string)
	Get(id string) (string, bool)
}

type Handler struct {
	storage storage.URLStorage
	baseURL string
}

func New(
	storage storage.URLStorage,
	baseURL string,
) *Handler {
	return &Handler{
		storage: storage,
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}

func (h *Handler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	// Нас интересует только POST /
	if r.Method != http.MethodPost || r.URL.Path != "/" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Проверяем Content-Type
	if !strings.HasPrefix(
		r.Header.Get("Content-Type"),
		"text/plain",
	) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Читаем всё тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	originalURL := strings.TrimSpace(string(body))

	if originalURL == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Генерируем уникальный ID
	var id string

	for {
		id = generateID()

		if _, exists := h.storage.Get(id); !exists {
			break
		}
	}

	h.storage.Save(id, originalURL)

	shortURL := fmt.Sprintf(
		"%s/%s",
		h.baseURL,
		id,
	)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)

	_, _ = w.Write([]byte(shortURL))
}

func (h *Handler) GetOriginalURL(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	originalURL, ok := h.storage.Get(id)

	if !ok {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func generateID() string {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	result := make([]byte, idLength)

	for i := range result {
		result[i] = alphabet[rand.Intn(len(alphabet))]
	}

	return string(result)
}
