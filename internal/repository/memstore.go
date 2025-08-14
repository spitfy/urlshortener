package repository

import (
	"context"
	"fmt"
	"sync"
)

type MemStore struct {
	mux *sync.Mutex
	s   map[string]string
}

func newMemStore() *MemStore {
	return &MemStore{
		mux: &sync.Mutex{},
		s:   make(map[string]string),
	}
}

func (s *MemStore) Add(_ context.Context, url URL, _ int) (hash string, err error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	if _, ok := s.s[url.Hash]; ok {
		return url.Hash, fmt.Errorf("wrong hash: '%s', already exists", url.Hash)
	}
	s.s[url.Hash] = url.Link
	return url.Hash, nil
}

func (s *MemStore) Get(_ context.Context, hash string) (string, error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	res, ok := s.s[hash]
	if !ok {
		return "", fmt.Errorf("data not found for n = %s", hash)
	}
	return res, nil
}

func (s *MemStore) Ping() error {
	return nil
}

func (s *MemStore) AllByUser(_ context.Context, _ int) ([]URL, error) {
	return make([]URL, 0), nil
}

func (s *MemStore) Close() error {
	return nil
}

func (s *MemStore) BatchAdd(ctx context.Context, urls []URL, userID int) error {
	for _, u := range urls {
		if _, err := s.Add(ctx, u, userID); err != nil {
			return err
		}
	}
	return nil
}

func (s *MemStore) CreateUser(_ context.Context) (int, error) {
	return -1, nil
}
