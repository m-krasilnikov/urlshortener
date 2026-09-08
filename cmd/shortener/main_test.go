package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/m-krasilnikov/urlshortener/internal/app"
	"github.com/m-krasilnikov/urlshortener/internal/storage"
)

func TestCreateShortURL(t *testing.T) {
	const baseURL = "http://localhost:8081"

	tests := []struct {
		name           string
		method         string
		path           string
		contentType    string
		body           string
		expectedStatus int
	}{
		{
			name:           "valid request",
			method:         http.MethodPost,
			path:           "/",
			contentType:    "text/plain",
			body:           "https://practicum.yandex.ru/",
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "wrong method",
			method:         http.MethodGet,
			path:           "/",
			contentType:    "text/plain",
			body:           "https://practicum.yandex.ru/",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "wrong path",
			method:         http.MethodPost,
			path:           "/test",
			contentType:    "text/plain",
			body:           "https://practicum.yandex.ru/",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "wrong content type",
			method:         http.MethodPost,
			path:           "/",
			contentType:    "application/json",
			body:           "https://practicum.yandex.ru/",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty body",
			method:         http.MethodPost,
			path:           "/",
			contentType:    "text/plain",
			body:           "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "text plain with charset",
			method:         http.MethodPost,
			path:           "/",
			contentType:    "text/plain; charset=utf-8",
			body:           "https://practicum.yandex.ru/",
			expectedStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := storage.NewMemoryStorage()
			router := app.NewRouter(st, baseURL)

			req := httptest.NewRequest(
				tt.method,
				tt.path,
				strings.NewReader(tt.body),
			)

			req.Header.Set(
				"Content-Type",
				tt.contentType,
			)

			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf(
					"expected status %d, got %d",
					tt.expectedStatus,
					rr.Code,
				)
			}
		})
	}
}

func TestCreateShortURLReturnsShortURL(t *testing.T) {
	const (
		baseURL     = "http://localhost:8081"
		originalURL = "https://practicum.yandex.ru/"
	)

	st := storage.NewMemoryStorage()
	router := app.NewRouter(st, baseURL)

	req := httptest.NewRequest(
		http.MethodPost,
		"/",
		strings.NewReader(originalURL),
	)

	req.Header.Set(
		"Content-Type",
		"text/plain",
	)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusCreated,
			rr.Code,
		)
	}

	contentType := rr.Header().Get("Content-Type")

	if !strings.HasPrefix(contentType, "text/plain") {
		t.Errorf(
			"expected Content-Type text/plain, got %q",
			contentType,
		)
	}

	shortURL := strings.TrimSpace(
		rr.Body.String(),
	)

	prefix := baseURL + "/"

	if !strings.HasPrefix(shortURL, prefix) {
		t.Fatalf(
			"expected short URL to start with %q, got %q",
			prefix,
			shortURL,
		)
	}

	id := strings.TrimPrefix(
		shortURL,
		prefix,
	)

	if len(id) != 8 {
		t.Errorf(
			"expected ID length %d, got %d",
			8,
			len(id),
		)
	}
}

func TestGetOriginalURL(t *testing.T) {
	const baseURL = "http://localhost:8081"

	st := storage.NewMemoryStorage()
	router := app.NewRouter(st, baseURL)

	// Сначала создаём короткую ссылку.
	originalURL := "https://practicum.yandex.ru/"

	createReq := httptest.NewRequest(
		http.MethodPost,
		"/",
		strings.NewReader(originalURL),
	)

	createReq.Header.Set(
		"Content-Type",
		"text/plain",
	)

	createRR := httptest.NewRecorder()

	router.ServeHTTP(
		createRR,
		createReq,
	)

	if createRR.Code != http.StatusCreated {
		t.Fatalf(
			"expected create status %d, got %d",
			http.StatusCreated,
			createRR.Code,
		)
	}

	shortURL := strings.TrimSpace(
		createRR.Body.String(),
	)

	id := strings.TrimPrefix(
		shortURL,
		baseURL+"/",
	)

	// Теперь проверяем GET /{id}.
	getReq := httptest.NewRequest(
		http.MethodGet,
		"/"+id,
		nil,
	)

	getRR := httptest.NewRecorder()

	router.ServeHTTP(
		getRR,
		getReq,
	)

	if getRR.Code != http.StatusTemporaryRedirect {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusTemporaryRedirect,
			getRR.Code,
		)
	}

	location := getRR.Header().Get("Location")

	if location != originalURL {
		t.Errorf(
			"expected Location %q, got %q",
			originalURL,
			location,
		)
	}
}

func TestGetOriginalURLUnknownID(t *testing.T) {

	st := storage.NewMemoryStorage()
	router := app.NewRouter(st,
		"http://localhost:8081",
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/Unknown1",
		nil,
	)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			rr.Code,
		)
	}
}

func TestInvalidRequests(t *testing.T) {
	st := storage.NewMemoryStorage()
	router := app.NewRouter(st,
		"http://localhost:8081",
	)

	tests := []struct {
		name   string
		method string
		path   string
	}{
		{
			name:   "unknown path",
			method: http.MethodGet,
			path:   "/test/test",
		},
		{
			name:   "POST with ID",
			method: http.MethodPost,
			path:   "/abc123",
		},
		{
			name:   "PUT root",
			method: http.MethodPut,
			path:   "/",
		},
		{
			name:   "DELETE root",
			method: http.MethodDelete,
			path:   "/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(
				tt.method,
				tt.path,
				nil,
			)

			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if rr.Code != http.StatusBadRequest {
				t.Errorf(
					"expected status %d, got %d",
					http.StatusBadRequest,
					rr.Code,
				)
			}
		})
	}
}
