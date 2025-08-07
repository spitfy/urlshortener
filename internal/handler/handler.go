package handler

import (
	"context"
	models "github.com/spitfy/urlshortener/internal/model"
	"net/http"
)

type Handler struct {
	service ServiceShortener
}

type ServiceShortener interface {
	Add(ctx context.Context, link string) (string, error)
	BatchAdd(ctx context.Context, req []models.BatchCreateRequest) ([]models.BatchCreateResponse, error)
	Get(ctx context.Context, hash string) (string, error)
	Ping() error
}

type RequestLogger interface {
	LogInfo(h http.HandlerFunc) http.HandlerFunc
}

var allowedContent = map[string]bool{
	"text/plain":         true,
	"application/json":   true,
	"application/x-gzip": true,
}

func newHandler(s ServiceShortener) *Handler {
	return &Handler{
		service: s,
	}
}
