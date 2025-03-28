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
	r.Use(middleware.WithAuth)

	// Add routes
	r.Post("/", h.PostURL)
	r.Post("/api/shorten/batch", h.ShortenBatchURL)
	r.Post("/api/shorten", h.ShortenURL)

	r.Get("/ping", h.IsAvailable)
	r.Get("/api/user/urls", h.GetUserURL)
	r.Get("/{urlKey}", h.GetURL)
	r.Get("/", h.GetURL)

	r.Delete("/api/user/urls", h.DeleteUserURL)

	return r
}
