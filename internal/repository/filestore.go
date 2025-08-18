package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spitfy/urlshortener/internal/config"
	models "github.com/spitfy/urlshortener/internal/model"
	"io"
	"os"
	"path/filepath"
)

type FileStore struct {
	file *os.File
	*MemStore
}

type LinkList []models.Link

func NewMockStore() *FileStore {
	return &FileStore{
		file:     nil,
		MemStore: newMemStore(),
	}
}

func newFileStore(config *config.Config) (*FileStore, error) {
	dir := filepath.Dir(config.FileStorage.FileStoragePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}
	f, err := os.OpenFile(config.FileStorage.FileStoragePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", config.FileStorage.FileStoragePath, err)
	}
	store := FileStore{
		file:     f,
		MemStore: newMemStore(),
	}
	links, err := store.init()
	if err != nil {
		return nil, fmt.Errorf("failed to init store: %w", err)
	}
	store.s = links

	return &store, nil
}

func (s *FileStore) Get(ctx context.Context, hash string) (URL, error) {
	return s.MemStore.Get(ctx, hash)
}

func (s *FileStore) getStore() (LinkList, error) {
	_, err := s.file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	var store LinkList
	dec := json.NewDecoder(s.file)
	err = dec.Decode(&store)
	if err == io.EOF {
		return store, nil
	}
	if err != nil {
		return nil, err
	}

	return store, nil
}

func (s *FileStore) init() (map[string]string, error) {
	_, err := s.file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	var store LinkList
	dec := json.NewDecoder(s.file)
	err = dec.Decode(&store)
	if err == io.EOF {
		store = nil
	} else if err != nil {
		return nil, err
	}
	links := make(map[string]string, len(store))
	for _, l := range store {
		links[l.ShortURL] = l.OriginalURL
	}
	return links, nil
}

func (s *FileStore) Add(ctx context.Context, url URL, userID int) (string, error) {
	if hash, err := s.MemStore.Add(ctx, url, userID); err != nil {
		return hash, err
	}
	if err := s.save(); err != nil {
		return url.Hash, err
	}
	return url.Hash, nil
}

func (s *FileStore) Ping() error {
	return nil
}

func (s *FileStore) Close() error {
	return nil
}

func (s *FileStore) BatchAdd(ctx context.Context, urls []URL, userID int) error {
	if err := s.MemStore.BatchAdd(ctx, urls, userID); err != nil {
		return err
	}
	if err := s.save(); err != nil {
		return err
	}
	return nil
}

func (s *FileStore) save() error {
	store := make(LinkList, 0, len(s.s))
	uuid := 1
	for hash, l := range s.s {
		ml := models.Link{
			UUID:        string(rune(uuid)),
			ShortURL:    hash,
			OriginalURL: l,
		}
		store = append(store, ml)
		uuid++
	}
	data, err := json.Marshal(store)
	if err != nil {
		return err
	}
	tmpPath := s.file.Name() + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return err
	}

	return os.Rename(tmpPath, s.file.Name())
}

func (s *FileStore) AllByUser(_ context.Context, _ int) ([]URL, error) {
	return make([]URL, 0), nil
}

func (s *FileStore) CreateUser(_ context.Context) (int, error) {
	return -1, nil
}
