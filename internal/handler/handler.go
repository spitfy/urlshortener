package handler

import (
	models "github.com/spitfy/urlshortener/internal/model"
	"net/http"
)

type Handler struct {
	service ServiceShortener
}

type ServiceShortener interface {
	Add(link string) (string, error)
	BatchAdd(req []models.BatchCreateRequest) ([]models.BatchCreateResponse, error)
	Get(hash string) (string, error)
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
