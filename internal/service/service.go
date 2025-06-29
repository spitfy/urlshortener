package service

import (
	"crypto/rand"
	"math/big"
	"net/url"

	"github.com/spitfy/urlshortener/internal/config"
	"github.com/spitfy/urlshortener/internal/repository"
)

const (
	chars   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	CharCnt = 8
)

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		b[i] = chars[num.Int64()]
	}
	return string(b)
}

type Service struct {
	store  repository.Store
	config config.Config
}

func (s *Service) Add(link string) string {
	hash := RandString(CharCnt)
	url := repository.URL{
		Link: link,
		Hash: hash,
	}
	s.store.Add(url)

	return s.makeURL(hash)
}

func (s *Service) Get(hash string) (string, error) {
	url, err := s.store.Get(hash)
	if err != nil {
		return "", err
	}

	return url, nil
}

func NewService(cfg config.Config) *Service {
	return &Service{
		store:  *repository.NewStore(),
		config: cfg,
	}
}

func NewMockService(cfg config.Config, r repository.Store) *Service {
	return &Service{
		store:  r,
		config: cfg,
	}
}

func (s *Service) makeURL(hash string) string {
	u := url.URL{
		Scheme: "http",
		Host:   s.config.Handlers.ServerAddr,
		Path:   "/",
	}

	return u.ResolveReference(&url.URL{Path: hash}).String()
}
