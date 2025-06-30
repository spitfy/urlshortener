package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/spitfy/urlshortener/internal/config"
	"github.com/spitfy/urlshortener/internal/repository"
	"github.com/spitfy/urlshortener/internal/service"
	"github.com/stretchr/testify/assert"
)

var (
	srv *httptest.Server
	cfg config.Config
)

func init() {
	cfg = config.GetConfig(false)
}

func TestMain(m *testing.M) {
	code := m.Run()
	srv.Close()
	os.Exit(code)
}

func TestHandler_Post(t *testing.T) {
	handler := newHandler(service.NewService(cfg))
	srv = httptest.NewServer(http.HandlerFunc(handler.Post))

	tests := []struct {
		name         string
		method       string
		expectedCode int
		url          string
	}{
		{
			name:         "success",
			method:       http.MethodPost,
			expectedCode: http.StatusCreated,
			url:          "https://pkg.go.dev/",
		},
		{
			name:         "wrong method",
			method:       http.MethodGet,
			expectedCode: http.StatusBadRequest,
			url:          "https://pkg.go.dev/",
		},
		{
			name:         "empty body",
			method:       http.MethodPost,
			expectedCode: http.StatusBadRequest,
			url:          "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := resty.New().R()
			req.SetHeader("Content-Type", "text/plain")

			if tt.url != "" {
				req.SetBody(tt.url)
			}

			resp, err := req.Execute(tt.method, srv.URL)

			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, tt.expectedCode, resp.StatusCode(), "Response code mismatch")

			if tt.expectedCode == http.StatusCreated {
				assert.Contains(t, string(resp.Body()), cfg.Handlers.ServerAddr)
			}
		})
	}
}

func TestHandler_Get(t *testing.T) {
	store := repository.NewStore()
	store.Add(repository.URL{Hash: "XXAABBOO", Link: "https://pkg.go.dev/"})

	handler := newHandler(service.NewMockService(cfg, *store))
	srv = httptest.NewServer(newRouter(handler))

	tests := []struct {
		name         string
		method       string
		expectedCode int
		hash         string
		location     string
	}{
		{
			name:         "success",
			method:       http.MethodGet,
			expectedCode: http.StatusOK,
			hash:         "XXAABBOO",
		},
		{
			name:         "not found",
			method:       http.MethodGet,
			expectedCode: http.StatusBadRequest,
			hash:         "UNKNOWN",
			location:     "",
		},
		{
			name:         "method not allowed",
			method:       http.MethodPost,
			expectedCode: http.StatusMethodNotAllowed,
			hash:         "XXAABBOO",
			location:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := resty.New()

			var resp *resty.Response
			var err error

			switch tt.method {
			case http.MethodGet:
				resp, err = client.R().
					SetHeader("Content-Type", "text/plain").
					Get(fmt.Sprintf("%s/%s", srv.URL, tt.hash))
			case http.MethodPost:
				resp, err = client.R().
					SetHeader("Content-Type", "text/plain").
					Post(fmt.Sprintf("%s/%s", srv.URL, tt.hash))
			default:
				t.Fatalf("unsupported method %s", tt.method)
			}

			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, tt.expectedCode, resp.StatusCode(), "Response code mismatch")
		})
	}
}
