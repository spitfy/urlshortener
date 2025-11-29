package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spitfy/urlshortener/internal/config"
	"github.com/spitfy/urlshortener/internal/model"
)

// FileStore реализует хранилище URL в файле с in-memory кэшем.
// Использует MemStore для быстрого доступа и синхронизирует данные с файлом.
// Пример создания:
//
//	conf := config.LoadConfig()
//	store, err := newFileStore(conf)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer store.Close()
type FileStore struct {
	file *os.File
	*MemStore
}

// LinkList представляет список ссылок в формате JSON файла
type LinkList []model.Link

// NewMockStore создает мок-хранилище без привязки к файлу (только in-memory).
// Пример:
//
//	mockStore := NewMockStore()
//	defer mockStore.Close()
func NewMockStore() *FileStore {
	return &FileStore{
		file:     nil,
		MemStore: newMemStore(),
	}
}

// newFileStore создает новое файловое хранилище с загрузкой данных из файла.
// Пример:
//
//	store, err := newFileStore(config)
//	if err != nil {
//	    // обработка ошибки
//	}
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

// GetByHash возвращает URL по хешу из in-memory кэша.
// Пример:
//
//	url, err := store.GetByHash(ctx, "abc123")
//	if err != nil {
//	    // обработка ошибки
//	}
func (s *FileStore) GetByHash(ctx context.Context, hash string) (URL, error) {
	return s.MemStore.GetByHash(ctx, hash)
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

// Add добавляет URL в хранилище с сохранением в файл.
// Пример:
//
//	hash, err := store.Add(ctx, URL{
//	    Hash: "abc",
//	    Link: "https://example.com"
//	}, 1)
//	if err != nil {
//	    // обработка ошибки
//	}
func (s *FileStore) Add(ctx context.Context, url URL, userID int) (string, error) {
	if hash, err := s.MemStore.Add(ctx, url, userID); err != nil {
		return hash, err
	}
	if err := s.save(); err != nil {
		return url.Hash, err
	}
	return url.Hash, nil
}

// Ping всегда возвращает nil (для совместимости с интерфейсом Storer).
// Пример:
//
//	if err := store.Ping(); err != nil {
//	    log.Println("Storage unavailable")
//	}
func (s *FileStore) Ping() error {
	return nil
}

// Close закрывает файловый дескриптор хранилища.
// Пример:
//
//	defer store.Close()
func (s *FileStore) Close() {}

// BatchAdd добавляет несколько URL с атомарным сохранением в файл.
// Пример:
//
//	urls := []URL{
//	    {Hash: "abc", Link: "https://example.com"},
//	    {Hash: "def", Link: "https://example.org"},
//	}
//	err := store.BatchAdd(ctx, urls, 1)
func (s *FileStore) BatchAdd(ctx context.Context, urls []URL, userID int) error {
	if err := s.MemStore.BatchAdd(ctx, urls, userID); err != nil {
		return err
	}
	if err := s.save(); err != nil {
		return err
	}
	return nil
}

// save сохраняет текущее состояние хранилища в файл
func (s *FileStore) save() error {
	store := make(LinkList, 0, len(s.s))
	uuid := 1
	for hash, l := range s.s {
		ml := model.Link{
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

// GetByUserID всегда возвращает пустой список (не реализовано для FileStore).
// Пример:
//
//	urls, _ := store.GetByUserID(ctx, 1) // всегда []
func (s *FileStore) GetByUserID(_ context.Context, _ int) ([]URL, error) {
	return make([]URL, 0), nil
}

// CreateUser всегда возвращает -1 (не реализовано для FileStore).
// Пример:
//
//	userID, _ := store.CreateUser(ctx) // всегда -1
func (s *FileStore) CreateUser(_ context.Context) (int, error) {
	return -1, nil
}

// BatchDelete всегда возвращает nil (не реализовано для FileStore).
// Пример:
//
//	_ = store.BatchDelete(ctx, UserHash{...})
func (s *FileStore) BatchDelete(_ context.Context, _ UserHash) (err error) {
	return nil
}

func (s *FileStore) Stats(ctx context.Context) (model.Stats, error) {
	return s.MemStore.Stats(ctx)
}
