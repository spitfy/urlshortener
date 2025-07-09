package repository

import (
	"fmt"
	"sync"
)

type Store struct {
	mux *sync.Mutex
	s   map[string]string
}

type URL struct {
	Link string
	Hash string
}

func NewStore() *Store {
	return &Store{
		mux: &sync.Mutex{},
		s:   make(map[string]string),
	}
}

func (s *Store) Add(url URL) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.s[url.Hash] = url.Link
}

func (s *Store) Get(hash string) (string, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	res, ok := s.s[hash]

	if !ok {
		return "", fmt.Errorf("data not found for n = %s", hash)
	}

	return res, nil
}
