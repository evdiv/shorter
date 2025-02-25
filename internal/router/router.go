package router

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"shorter/internal/handlers"
	"shorter/internal/middleware"
)

func NewRouter(h *handlers.Handlers) http.Handler {
	r := chi.NewRouter()

	// Add middleware
	r.Use(middleware.WithLogging)
	r.Use(middleware.WithGzip)

	// Add routes
	r.Post("/", h.PostURL)
	r.Post("/api/shorten", h.ShortenURL)
	r.Get("/{urlKey}", h.GetURL)
	r.Get("/", h.GetURL)

	return r
}
