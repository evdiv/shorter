package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"shorter/cmd/shortener/config"
	"strings"
	"testing"
)

func setupRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/", PostURL)
	r.Get("/{urlKey}", GetURL)
	r.Get("/", GetURL)
	return r
}

func TestRouter(t *testing.T) {
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
			body:   "https://practicum.yandex.ru/",
			want: want{
				code:   201,
				header: "",
				body:   config.GetHost("Result") + "/9740", // Expected body will vary depending on generated key
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
				body:   "The body should contain URL",
			},
		},
		{
			name:   "GET: Positive. Extract URL by a valid key",
			target: "/9740", // Use the expected key from the POST test
			method: "GET",
			body:   "",
			want: want{
				code:   307,
				header: "https://practicum.yandex.ru/",
				body:   "",
			},
		},
		{
			name:   "GET: Negative. Extract URL by non-valid key",
			target: "/xxxxxxxxx",
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
