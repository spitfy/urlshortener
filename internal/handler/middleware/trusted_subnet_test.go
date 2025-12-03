package middleware

import (
	"github.com/spitfy/urlshortener/internal/config"
	config2 "github.com/spitfy/urlshortener/internal/handler/config"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTrustedSubnet(t *testing.T) {
	cfg := config.Config{Handlers: config2.Config{TrustedSubnet: "192.168.1.0/24"}}
	mw := TrustedSubnet(&cfg)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Real-IP", "192.168.1.100")
	rr := httptest.NewRecorder()
	mw(func(w http.ResponseWriter, r *http.Request) {}).ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)

	req.Header.Set("X-Real-IP", "10.0.0.1")
	rr = httptest.NewRecorder()
	mw(func(w http.ResponseWriter, r *http.Request) {}).ServeHTTP(rr, req)
	require.Equal(t, http.StatusForbidden, rr.Code)
}
