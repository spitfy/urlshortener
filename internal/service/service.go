package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/spitfy/urlshortener/internal/audit"
	models "github.com/spitfy/urlshortener/internal/model"
	"log"
	"math/big"
	"net/url"
	"runtime"
	"sync"

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
	store     repository.Storer
	config    config.Config
	deleteQ   chan repository.UserHash
	observers []audit.Observer
	mu        sync.Mutex
}

func NewService(cfg config.Config, store repository.Storer) *Service {
	s := &Service{
		store:   store,
		config:  cfg,
		deleteQ: make(chan repository.UserHash, 100),
	}

	maxProcs := runtime.GOMAXPROCS(0)
	for i := 0; i < maxProcs; i++ {
		go s.runDeleteWorker()
	}

	return s
}

func (s *Service) runDeleteWorker() {
	for uh := range s.deleteQ {
		if err := s.store.BatchDelete(context.Background(), uh); err != nil {
			log.Printf("batch delete error: %v", err)
		}
	}
}

func (s *Service) DeleteEnqueue(_ context.Context, hashes []string, userID int) {
	s.deleteQ <- repository.UserHash{
		UserID: userID,
		Hash:   hashes,
	}
}

func (s *Service) Add(ctx context.Context, link string, userID int) (string, error) {
	if !isURL(link) {
		return "", errors.New("invalid url")
	}

	hash := RandString(CharCnt)
	u := repository.URL{Link: link, Hash: hash}
	hash, err := s.store.Add(ctx, u, userID)

	if err != nil && !errors.Is(err, repository.ErrExistsURL) {
		return "", err
	}
	shortURL, errMakeURL := s.makeURL(hash)
	if errMakeURL != nil {
		return "", errMakeURL
	}
	if errors.Is(err, repository.ErrExistsURL) {
		return shortURL, repository.ErrExistsURL
	}
	return shortURL, nil
}

func (s *Service) BatchAdd(
	ctx context.Context,
	req []models.BatchCreateRequest,
	userID int,
) ([]models.BatchCreateResponse, error) {
	res := make([]models.BatchCreateResponse, 0, len(req))
	for _, r := range req {
		shortURL, err := s.Add(ctx, r.OriginalURL, userID)
		if err != nil {
			return nil, err
		}
		res = append(res, models.BatchCreateResponse{CorrelationID: r.CorrelationID, ShortURL: shortURL})
	}
	return res, nil
}

func (s *Service) GetByHash(ctx context.Context, hash string) (repository.URL, error) {
	return s.store.GetByHash(ctx, hash)
}

func (s *Service) Ping() error {
	return s.store.Ping()
}

func (s *Service) GetByUserID(ctx context.Context, id int) ([]models.LinkPair, error) {
	links, err := s.store.GetByUserID(ctx, id)
	if err != nil {
		return nil, err
	}
	res := make([]models.LinkPair, 0, len(links))
	for _, u := range links {
		ShortURL, err := s.makeURL(u.Hash)
		if err != nil {
			return nil, err
		}
		res = append(res, models.LinkPair{
			ShortURL:    ShortURL,
			OriginalURL: u.Link,
		})
	}
	return res, nil
}

func (s *Service) CreateUser(ctx context.Context) (int, error) {
	id, err := s.store.CreateUser(ctx)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Service) makeURL(hash string) (string, error) {
	addr, err := url.JoinPath(s.config.Service.ServerURL, hash)
	if err != nil {
		return "", fmt.Errorf("can't create short url: %w", err)
	}
	return addr, nil
}

func isURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func (s *Service) AddObserver(observer audit.Observer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.observers = append(s.observers, observer)
}

func (s *Service) NotifyObservers(ctx context.Context, event audit.Event) {
	s.mu.Lock()
	observers := make([]audit.Observer, len(s.observers))
	copy(observers, s.observers)
	s.mu.Unlock()

	for _, observer := range observers {
		go observer.Notify(ctx, event)
	}
}
