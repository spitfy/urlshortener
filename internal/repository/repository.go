package repository

import (
	"errors"
	"fmt"
	"sync"
)

type Store struct {
	mux *sync.Mutex
	s   map[string]string
}

type Url struct {
	Link string
	Hash string
}

var (
	ErrGetUrlNotFound = errors.New("data not found")
)

func newErrGetUrlNotFound(hash string) error {
	return fmt.Errorf("%w for n = %s", ErrGetUrlNotFound, hash)
}

func NewStore() *Store {
	return &Store{
		mux: &sync.Mutex{},
		s:   make(map[string]string),
	}
}

func (s *Store) Add(url Url) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.s[url.Hash] = url.Link
}

func (s *Store) Get(hash string) (string, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	res, ok := s.s[hash]

	if !ok {
		return "", newErrGetUrlNotFound(hash)
	}

	return res, nil
}
