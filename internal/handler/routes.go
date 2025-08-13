package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/spitfy/urlshortener/internal/auth"
	"github.com/spitfy/urlshortener/internal/config"
)

func Serve(cfg config.Config, service ServiceShortener, l RequestLogger) error {
	a := auth.New(cfg.Auth.SecretKey)
	h := newHandler(service, a)
	router := newRouter(h, l)

	server := &http.Server{
		Addr:    cfg.Handlers.ServerAddr,
		Handler: router,
	}

	return server.ListenAndServe()
}

func newRouter(h *Handler, l RequestLogger) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/ping", h.authMiddleware(gzipMiddleware(l.LogInfo(h.Ping))))
	r.Get("/{hash}", h.authMiddleware(gzipMiddleware(l.LogInfo(h.Get))))
	r.Get("/api/user/urls", h.authMiddleware(gzipMiddleware(l.LogInfo(h.GetByUserId))))
	r.Post("/api/shorten/batch", h.authMiddleware(gzipMiddleware(l.LogInfo(h.Batch))))
	r.Post("/api/shorten", h.authMiddleware(h.ShortenURL))
	r.Post("/", h.authMiddleware(gzipMiddleware(l.LogInfo(h.Post))))

	return r
}
