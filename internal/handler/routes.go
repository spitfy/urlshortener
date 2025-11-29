// Package handler содержит маршрутизацию HTTP-запросов.
package handler

import (
	"fmt"
	"github.com/spitfy/urlshortener/internal/gomodule"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"net/http/pprof"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/spitfy/urlshortener/internal/auth"
	"github.com/spitfy/urlshortener/internal/config"
)

// Serve запускает HTTP-сервер с обработчиками URL shortener API.
// Принимает конфигурацию, сервис сокращения URL и логгер запросов.
// Возвращает сервер и ошибку в случае неудачного запуска сервера.
func Serve(cfg config.Config, service ServiceShortener, l RequestLogger) (*http.Server, error) {
	a := auth.New(cfg.Auth.SecretKey)
	h := newHandler(service, a)
	router := newRouter(h, l)

	server := &http.Server{
		Addr:         cfg.Handlers.ServerAddr,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if cfg.Handlers.EnableHTTPS {
		certFile, err := CertPath(cfg.Handlers.CertFile)
		if err != nil {
			return server, fmt.Errorf("certificate file not readable: %s — %w", certFile, err)
		}
		keyFile, err := CertPath(cfg.Handlers.KeyFile)
		if err != nil {
			return server, fmt.Errorf("key file not readable: %s — %w", keyFile, err)
		}
		httpsAddr := ":" + cfg.Handlers.HTTPSPort
		server.Addr = httpsAddr
		fmt.Printf("Starting HTTPS server on %s\n", httpsAddr)
	}

	return server, nil
}

// newRouter создает новый маршрутизатор с обработчиками для:
// - API сокращения URL
// - Профилирования (pprof)
// Добавляет middleware для аутентификации, сжатия и логирования.
func newRouter(h *Handler, l RequestLogger) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		wd, _ := os.Getwd()
		root, _ := gomodule.FindModuleRoot(wd)
		f := filepath.Join(root, "docs", "swagger.json")
		http.ServeFile(w, r, f)
	})
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger.json"),
	))

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

func CertPath(cert string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	moduleRoot, err := gomodule.FindModuleRoot(wd)
	if err != nil {
		return "", err
	}

	certFile := filepath.Join(moduleRoot, cert)
	if _, err = os.Stat(certFile); err != nil {
		return "", err
	}

	return certFile, nil
}
