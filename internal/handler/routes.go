package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/spitfy/urlshortener/internal/config"
)

func Serve(cfg config.Config, service ServiceShortener, l RequestLogger) error {
	h := newHandler(service)
	router := newRouter(h, l)

	server := &http.Server{
		Addr:    cfg.Handlers.ServerAddr,
		Handler: router,
	}

	return server.ListenAndServe()
}

func newRouter(h *Handler, l RequestLogger) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/ping", gzipMiddleware(l.LogInfo(h.Ping)))
	r.Get("/{hash}", gzipMiddleware(l.LogInfo(h.Get)))
	r.Post("/api/shorten/batch", gzipMiddleware(l.LogInfo(h.Batch)))
	r.Post("/api/shorten", gzipMiddleware(l.LogInfo(h.ShortenURL)))
	r.Post("/", gzipMiddleware(l.LogInfo(h.Post)))

	return r
}
