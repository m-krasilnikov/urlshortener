package handler

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/m-krasilnikov/urlshortener/internal/storage"
)

const idLength = 8

type Handler struct {
	storage *storage.MemoryStorage
	baseURL string
}

func New(
	storage *storage.MemoryStorage,
	baseURL string,
) *Handler {
	return &Handler{
		storage: storage,
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}

func (h *Handler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	// Нас интересует только POST /
	log.Println(">>> CreateShortURL ENTER")

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
	log.Println(">>> GetOriginalURL ENTER")
	// Нас интересует только GET
	if r.Method != http.MethodGet {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/")

	// ID должен присутствовать
	if id == "" || strings.Contains(id, "/") {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	originalURL, ok := h.storage.Get(id)

	if !ok {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// 307 Temporary Redirect
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
