package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"shorter/internal/config"
	"shorter/internal/storage"
	"strings"
	"testing"
)

func setupRouter() *chi.Mux {

	memStorage := storage.NewMemoryStorage()
	h := NewHandlers(memStorage)

	r := chi.NewRouter()
	r.Post("/", h.PostURL)
	r.Post("/api/shorten", h.ShortenURL)
	r.Get("/{urlKey}", h.GetURL)
	r.Get("/", h.GetURL)

	return r
}

func TestRouter(t *testing.T) {

	resultHost := config.GetHost("Result")

	type want struct {
		code   int
		header string
		body   string
	}

	tests := []struct {
		name   string
		target string
		method string
		body   string
		want   want
	}{
		{
			name:   "POST: Positive. Valid URL for storing",
			target: "/",
			method: "POST",
			body:   "https://practicum.yandex.ru",
			want: want{
				code:   201,
				header: "",
				body:   resultHost + "/921c", // Expected body will vary depending on generated key
			},
		},
		{
			name:   "POST: Negative. Empty URL for storing",
			target: "/",
			method: "POST",
			body:   "",
			want: want{
				code:   400,
				header: "",
				body:   "The body should contain a valid URL",
			},
		},
		{
			name:   "GET: Positive. Extract URL by a valid key",
			target: "/921c", // Use the expected key from the POST test
			method: "GET",
			body:   "",
			want: want{
				code:   307,
				header: "https://practicum.yandex.ru",
				body:   "",
			},
		},
		{
			name:   "GET: Negative. Extract URL by non-valid key",
			target: "/error",
			method: "GET",
			body:   "",
			want: want{
				code:   400,
				header: "",
				body:   "URL is not found",
			},
		},
		{
			name:   "GET: Negative. Extract URL by empty key",
			target: "/",
			method: "GET",
			body:   "",
			want: want{
				code:   400,
				header: "",
				body:   "Missing URL key",
			},
		},
		{
			name:   "POST: Positive. JSON formatted URL for storing",
			target: "/api/shorten",
			method: "POST",
			body:   `{"url":"https://practicum.yandex.ru"}`,
			want: want{
				code:   201,
				header: "",
				body:   `{"result":"` + resultHost + `/921c"}`, // Expected body will vary depending on generated key
			},
		},
	}

	router := setupRouter()

	// Run the tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.target, strings.NewReader(tt.body))

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.code, res.StatusCode)

			if tt.want.header != "" {
				assert.Equal(t, tt.want.header, res.Header.Get("Location"))
			}
			body := strings.TrimSpace(w.Body.String())
			assert.Equal(t, tt.want.body, body)

		})
	}
}
