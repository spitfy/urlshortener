package handler

import (
	"context"
	"github.com/spitfy/urlshortener/internal/auth"
	models "github.com/spitfy/urlshortener/internal/model"
	"github.com/spitfy/urlshortener/internal/repository"
	"net/http"
)

type Handler struct {
	service ServiceShortener
	auth    *auth.AuthManager
}

type ServiceShortener interface {
	Add(ctx context.Context, link string, userID int) (string, error)
	BatchAdd(ctx context.Context, req []models.BatchCreateRequest, userID int) ([]models.BatchCreateResponse, error)
	Get(ctx context.Context, hash string) (repository.URL, error)
	Ping() error
	GetByUserID(ctx context.Context, userID int) ([]models.LinkPair, error)
	CreateUser(ctx context.Context) (int, error)
}

type RequestLogger interface {
	LogInfo(h http.HandlerFunc) http.HandlerFunc
}

var allowedContent = map[string]bool{
	"text/plain":         true,
	"application/json":   true,
	"application/x-gzip": true,
}

func newHandler(s ServiceShortener, a *auth.AuthManager) *Handler {
	return &Handler{
		service: s,
		auth:    a,
	}
}
