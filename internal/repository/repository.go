package repository

import (
	"encoding/json"
	"fmt"
	"github.com/spitfy/urlshortener/internal/config"
	models "github.com/spitfy/urlshortener/internal/model"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type Store struct {
	mux  *sync.RWMutex
	s    map[string]link
	file *os.File
}

type URL struct {
	Link string
	Hash string
}

type link struct {
	URL  string
	UUID string
}

func NewMockStore() *Store {
	return &Store{
		mux:  &sync.RWMutex{},
		s:    nil,
		file: nil,
	}
}

func NewStore(config *config.Config) *Store {
	fullPath := filepath.Join(config.FileStorage.FileStoragePath, config.FileStorage.FileStorageName)
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(fmt.Errorf("failed to create storage directory: %w", err))
	}
	f, err := os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	store := Store{
		mux:  &sync.RWMutex{},
		s:    nil,
		file: f,
	}
	links, err := store.init()
	if err != nil {
		panic(err)
	}
	store.s = *links

	return &store
}

func (s *Store) Add(url URL) error {
	s.mux.Lock()
	defer s.mux.Unlock()
	uuid := len(s.s) + 1
	s.s[url.Hash] = link{url.Link, fmt.Sprintf("%d", uuid)}

	return s.save()
}

func (s *Store) Get(hash string) (string, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	res, ok := s.s[hash]
	if !ok {
		return "", fmt.Errorf("data not found for n = %s", hash)
	}

	return res.URL, nil
}

func (s *Store) getStore() (*models.Store, error) {
	_, err := s.file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	var store models.Store
	dec := json.NewDecoder(s.file)
	err = dec.Decode(&store)
	if err == io.EOF {
		return &store, nil
	}
	if err != nil {
		return nil, err
	}

	return &store, nil
}

func (s *Store) init() (*map[string]link, error) {
	store, err := s.getStore()
	if err != nil {
		return nil, err
	}
	links := make(map[string]link, len(*store))
	for _, l := range *store {
		links[l.ShortURL] = link{l.OriginalURL, l.UUID}
	}
	return &links, nil
}

func (s *Store) save() error {
	var store models.Store
	for hash, l := range s.s {
		ml := models.Link{
			UUID:        l.UUID,
			ShortURL:    hash,
			OriginalURL: l.URL,
		}
		store = append(store, ml)
	}

	data, _ := json.Marshal(store)
	tmpPath := s.file.Name() + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return err
	}

	return os.Rename(tmpPath, s.file.Name())
}
