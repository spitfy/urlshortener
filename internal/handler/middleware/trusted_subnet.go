package middleware

import (
	"github.com/spitfy/urlshortener/internal/config"
	"log"
	"net"
	"net/http"
	"strings"
)

// TrustedSubnet создает middleware для проверки IP
func TrustedSubnet(cfg *config.Config) func(next http.HandlerFunc) http.HandlerFunc {
	var (
		trustedNet *net.IPNet
		err        error
	)

	if cfg.Handlers.TrustedSubnet != "" {
		_, trustedNet, err = net.ParseCIDR(cfg.Handlers.TrustedSubnet)
		if err != nil {
			log.Printf("invalid trusted_subnet %s: %v", cfg.Handlers.TrustedSubnet, err)
		}
	}

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if cfg.Handlers.TrustedSubnet == "" {
				http.Error(w, "Access forbidden", http.StatusForbidden)
				return
			}

			clientIP := strings.TrimSpace(r.Header.Get("X-Real-IP"))
			if clientIP == "" {
				http.Error(w, "X-Real-IP header required", http.StatusForbidden)
				return
			}

			ip := net.ParseIP(clientIP)
			if ip == nil {
				http.Error(w, "Invalid X-Real-IP format", http.StatusForbidden)
				return
			}

			// Проверяем вхождение в подсеть
			if !trustedNet.Contains(ip) {
				http.Error(w, "IP not in trusted subnet", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}
