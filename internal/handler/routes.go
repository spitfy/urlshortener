package handler

import (
	"net/http"
	"net/http/pprof"

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
	r.Get("/api/user/urls", h.authMiddleware(gzipMiddleware(l.LogInfo(h.GetByUserID))))
	r.Delete("/api/user/urls", h.authMiddleware(gzipMiddleware(l.LogInfo(h.Delete))))
	r.Post("/api/shorten/batch", h.authMiddleware(gzipMiddleware(l.LogInfo(h.BatchAdd))))
	r.Post("/api/shorten", h.authMiddleware(gzipMiddleware(l.LogInfo(h.ShortenURL))))
	r.Post("/", h.authMiddleware(gzipMiddleware(l.LogInfo(h.Post))))

	r.Group(func(r chi.Router) {
		r.Handle("/debug/pprof/*", http.HandlerFunc(pprof.Index))
		r.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
		r.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
		r.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
		r.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	})
	return r
}
