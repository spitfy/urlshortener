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

func (s *MemStore) Add(url URL) error {
	s.mux.Lock()
	defer s.mux.Unlock()
	if _, ok := s.s[url.Hash]; ok {
		return fmt.Errorf("wrong hash: '%s', already exists", url.Hash)
	}
	s.s[url.Hash] = url.Link
	return nil
}

func (s *MemStore) Get(hash string) (string, error) {
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

func (s *MemStore) Close() error {
	return nil
}

func (s *MemStore) BatchAdd(_ context.Context, urls []URL) error {
	for _, u := range urls {
		if err := s.Add(u); err != nil {
			return err
		}
	}
	return nil
}
