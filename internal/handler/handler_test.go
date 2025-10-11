package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spitfy/urlshortener/internal/auth"
	authConf "github.com/spitfy/urlshortener/internal/auth/config"
	handlerConf "github.com/spitfy/urlshortener/internal/handler/config"
	"github.com/spitfy/urlshortener/internal/logger"
	models "github.com/spitfy/urlshortener/internal/model"
	repoConf "github.com/spitfy/urlshortener/internal/repository/config"
	serviceConf "github.com/spitfy/urlshortener/internal/service/config"
	"github.com/stretchr/testify/require"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/spitfy/urlshortener/internal/config"
	"github.com/spitfy/urlshortener/internal/repository"
	"github.com/spitfy/urlshortener/internal/service"
	"github.com/stretchr/testify/assert"
)

var (
	srv *httptest.Server
	cfg = config.Config{
		Handlers:    handlerConf.Config{ServerAddr: config.DefaultServerAddr},
		Service:     serviceConf.Config{ServerURL: config.DefaultServerURL},
		FileStorage: repoConf.Config{FileStoragePath: config.DefaultFileStorageTest},
		Auth:        authConf.Config{SecretKey: config.SecretKey},
	}
	am = auth.New(cfg.Auth.SecretKey)
)

func TestMain(m *testing.M) {
	code := m.Run()
	if srv != nil {
		srv.Close()
	}
	if err := os.Remove(cfg.FileStorage.FileStoragePath); err != nil {
		log.Println(err)
	}
	os.Exit(code)
}

func TestHandler_Post(t *testing.T) {
	store, err := repository.CreateStore(&cfg)
	require.NoError(t, err, "error creating store")
	h := newHandler(service.NewService(cfg, store), am)
	srv = httptest.NewServer(h.authMiddleware(h.Post))

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
	store, err := repository.CreateStore(&cfg)
	require.NoError(t, err, "error creating store")
	ctx := context.Background()
	_, _ = store.Add(ctx, repository.URL{Hash: "XXAABBOO", Link: "https://pkg.go.dev/"}, -1)
	handler := newHandler(service.NewService(cfg, store), am)
	l := logger.InitMock()
	srv = httptest.NewServer(newRouter(handler, l))

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
			expectedCode: http.StatusTemporaryRedirect,
			hash:         "XXAABBOO",
			location:     "https://pkg.go.dev/",
		},
		{
			name:         "not_found",
			method:       http.MethodGet,
			expectedCode: http.StatusBadRequest,
			hash:         "UNKNOWN",
			location:     "",
		},
		{
			name:         "method_not_allowed",
			method:       http.MethodPost,
			expectedCode: http.StatusMethodNotAllowed,
			hash:         "XXAABBOO",
			location:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := resty.New()
			client.SetRedirectPolicy(resty.NoRedirectPolicy())

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
				t.Errorf("unsupported method %s", tt.method)
				return
			}

			if err != nil && !strings.Contains(err.Error(), "auto redirect is disabled") {
				assert.NoError(t, err, "error making HTTP request")
			}

			if tt.location != "" {
				assert.Equal(t, tt.location, resp.Header().Get("Location"))
			}

			assert.Equal(t, tt.expectedCode, resp.StatusCode(), "Response code mismatch")
		})
	}
}

func TestHandler_ShortenURL(t *testing.T) {
	store, err := repository.CreateStore(&cfg)
	require.NoError(t, err, "error creating store")
	h := newHandler(service.NewService(cfg, store), am)
	srv = httptest.NewServer(h.authMiddleware(h.ShortenURL))
	token, _ := am.BuildJWT(123)

	tests := []struct {
		name         string
		method       string
		body         string
		contentType  string
		expectedCode int
		expectedBody bool
	}{
		{
			name:         "method_get",
			method:       http.MethodGet,
			body:         "",
			contentType:  "application/json",
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: false,
		},
		{
			name:         "method_put",
			method:       http.MethodGet,
			body:         "",
			contentType:  "application/json",
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: false,
		},
		{
			name:         "method_post_without_body",
			method:       http.MethodPost,
			body:         "",
			contentType:  "application/json",
			expectedCode: http.StatusBadRequest,
			expectedBody: false,
		},
		{
			name:         "method_get_unsupported_type",
			method:       http.MethodGet,
			body:         `{"res": "https://www.perplexity.ai"}`,
			contentType:  "application/json",
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: false,
		},
		{
			name:         "method_post_success",
			method:       http.MethodPost,
			body:         `{"url": "https://www.perplexity.ai"}`,
			contentType:  "application/json",
			expectedCode: http.StatusCreated,
			expectedBody: true,
		},
		{
			name:         "bad_content_type",
			method:       http.MethodPost,
			body:         `{"url": "https://www.perplexity.ai"}`,
			contentType:  "text/plain",
			expectedCode: http.StatusBadRequest,
			expectedBody: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := resty.New()
			//client.SetDebug(true)
			resp, err := client.R().
				SetHeader("Content-Type", tt.contentType).
				SetBody(tt.body).
				SetCookie(&http.Cookie{
					Name:  "ID",
					Value: token,
					Path:  "/",
				}).
				Execute(tt.method, srv.URL)

			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, tt.expectedCode, resp.StatusCode(), "Response code mismatch")

			if tt.expectedBody {
				var respBody models.Response
				err := json.Unmarshal(resp.Body(), &respBody)
				assert.NoError(t, err, "response is not valid JSON")
				assert.Contains(t, respBody.Result, cfg.Handlers.ServerAddr, "result must start with base url")
			}
		})
	}
}
