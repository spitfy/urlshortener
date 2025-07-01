package service

import (
	"crypto/rand"
	"errors"
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
	store  Storer
	config config.Config
}

type Storer interface {
	Add(url repository.URL)
	Get(hash string) (string, error)
}

func (s *Service) Add(link string) (string, error) {
	hash := RandString(CharCnt)
	url := repository.URL{
		Link: link,
		Hash: hash,
	}
	s.store.Add(url)

	link, err := s.makeURL(hash)
	if err != nil {
		return "", err
	}

	return link, nil
}

func (s *Service) Get(hash string) (string, error) {
	url, err := s.store.Get(hash)
	if err != nil {
		return "", err
	}

	return url, nil
}

func NewService(cfg config.Config, store Storer) *Service {
	return &Service{
		store:  store,
		config: cfg,
	}
}

func (s *Service) makeURL(hash string) (string, error) {
	addr, err := url.JoinPath(s.config.Service.ServerURL, hash)

	if err != nil {
		return "", errors.New("can't create short url")
	}

	return addr, nil
}
