package repository

import (
	"context"
	"fmt"
	"github.com/spitfy/urlshortener/internal/model"
	"sync"
)

// MemStore реализует хранилище URL в памяти с синхронизацией доступа.
// Пример создания:
//
//	store := newMemStore()
type MemStore struct {
	mux *sync.Mutex
	s   map[string]string
}

// newMemStore создает новый экземпляр MemStore.
// Пример:
//
//	store := newMemStore()
func newMemStore() *MemStore {
	return &MemStore{
		mux: &sync.Mutex{},
		s:   make(map[string]string),
	}
}

// Add добавляет URL в хранилище.
// Возвращает ошибку если хеш уже существует.
// Пример:
//
//	hash, err := store.Add(ctx, URL{Hash: "abc", Link: "https://example.com"}, 1)
//	if err != nil {
//	    // обработка ошибки
//	}
func (s *MemStore) Add(_ context.Context, url URL, _ int) (hash string, err error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	if _, ok := s.s[url.Hash]; ok {
		return url.Hash, fmt.Errorf("wrong hash: '%s', already exists", url.Hash)
	}
	s.s[url.Hash] = url.Link
	return url.Hash, nil
}

// GetByHash возвращает URL по его хешу.
// Возвращает ошибку если URL не найден.
// Пример:
//
//	url, err := store.GetByHash(ctx, "abc")
//	if err != nil {
//	    // обработка ошибки
//	}
func (s *MemStore) GetByHash(_ context.Context, hash string) (URL, error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	link, ok := s.s[hash]
	if !ok {
		return URL{}, fmt.Errorf("data not found for n = %s", hash)
	}
	return URL{
		Link: link,
		Hash: hash,
	}, nil
}

// Ping всегда возвращает nil (для совместимости с интерфейсом Storer).
// Пример:
//
//	if err := store.Ping(); err != nil {
//	    log.Println("Ошибка соединения с хранилищем")
//	}
func (s *MemStore) Ping() error {
	return nil
}

// GetByUserID возвращает пустой список URL (in-memory хранилище не поддерживает пользовательские URL).
// Пример:
//
//	links, _ := store.GetByUserID(ctx, 1) // всегда возвращает пустой slice
func (s *MemStore) GetByUserID(_ context.Context, _ int) ([]URL, error) {
	return make([]URL, 0), nil
}

func (s *MemStore) Close() {}

// BatchAdd добавляет несколько URL в хранилище атомарно.
// Возвращает первую ошибку, если добавление не удалось.
// Пример:
//
//	urls := []URL{
//	    {Hash: "abc", Link: "https://example.com"},
//	    {Hash: "def", Link: "https://example.org"},
//	}
//	err := store.BatchAdd(ctx, urls, 1)
func (s *MemStore) BatchAdd(ctx context.Context, urls []URL, userID int) error {
	for _, u := range urls {
		if _, err := s.Add(ctx, u, userID); err != nil {
			return err
		}
	}
	return nil
}

// CreateUser всегда возвращает -1 (in-memory хранилище не поддерживает пользователей).
// Пример:
//
//	userID, _ := store.CreateUser(ctx) // всегда возвращает -1
func (s *MemStore) CreateUser(_ context.Context) (int, error) {
	return -1, nil
}

// BatchDelete всегда возвращает nil (in-memory хранилище не поддерживает удаление).
// Пример:
//
//	err := store.BatchDelete(ctx, UserHash{UserID: 1, Hash: []string{"abc"}})
func (s *MemStore) BatchDelete(_ context.Context, _ UserHash) (err error) {
	return nil
}

// Stats статистика по количеству ссылок в сервисе
func (s *MemStore) Stats(_ context.Context) (model.Stats, error) {
	return model.Stats{URLs: len(s.s), Users: 1}, nil
}
