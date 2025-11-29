// Package handler содержит основные HTTP-обработчики приложения.
package handler

import (
	"context"
	"github.com/spitfy/urlshortener/internal/audit"
	"github.com/spitfy/urlshortener/internal/auth"
	"github.com/spitfy/urlshortener/internal/model"
	"github.com/spitfy/urlshortener/internal/repository"
	"net/http"
)

type Handler struct {
	service ServiceShortener
	auth    auth.AuthManager
}

type ServiceShortener interface {
	Add(ctx context.Context, link string, userID int) (string, error)
	BatchAdd(ctx context.Context, req []model.BatchCreateRequest, userID int) ([]model.BatchCreateResponse, error)
	GetByHash(ctx context.Context, hash string) (repository.URL, error)
	Ping() error
	GetByUserID(ctx context.Context, userID int) ([]model.LinkPair, error)
	CreateUser(ctx context.Context) (int, error)
	DeleteEnqueue(ctx context.Context, req []string, userID int)
	AddObserver(observer audit.Observer)
	NotifyObservers(ctx context.Context, event audit.Event)
	Stats(ctx context.Context) (model.Stats, error)
}

type RequestLogger interface {
	LogInfo(h http.HandlerFunc) http.HandlerFunc
}

var allowedContent = map[string]bool{
	"text/plain":         true,
	"application/json":   true,
	"application/x-gzip": true,
}

func newHandler(s ServiceShortener, a *auth.Manager) *Handler {
	return &Handler{
		service: s,
		auth:    a,
	}
}
