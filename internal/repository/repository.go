package repository

import (
	"context"
	"errors"
	"github.com/spitfy/urlshortener/internal/model"

	"github.com/spitfy/urlshortener/internal/config"
)

// URL представляет структуру для хранения информации о сокращенной ссылке.
// Пример:
//
//	url := URL{
//	    Link: "https://example.com",
//	    Hash: "abc123",
//	    DeletedFlag: false,
//	}
type URL struct {
	Link        string // Оригинальный URL
	Hash        string // Сокращенный идентификатор
	DeletedFlag bool   // Флаг удаления (soft delete)
}

// UserHash содержит информацию о пользователе и хешах для пакетных операций.
// Пример:
//
//	userHash := UserHash{
//	    UserID: 1,
//	    Hash:   []string{"abc123", "def456"},
//	}
type UserHash struct {
	UserID int      // Идентификатор пользователя
	Hash   []string // Список хешей для операций
}

// ErrExistsURL возвращается при попытке добавить уже существующий URL.
var ErrExistsURL = errors.New("URL already exists")

// Storer определяет интерфейс для работы с хранилищем URL.
// Реализации:
//   - DBStore (PostgreSQL)
//   - FileStore (файловое хранилище)
//   - MemStore (in-memory)
//
//go:generate mockgen -destination=storer_mock.go -package=order github.com/spitfy/urlshortener/internal/repository Storer
type Storer interface {
	// Add добавляет новую ссылку в хранилище.
	// Возвращает ErrExistsURL если URL уже существует.
	// Пример:
	//   hash, err := store.Add(ctx, URL{...}, userID)
	Add(ctx context.Context, url URL, userID int) (hash string, err error)

	// GetByHash возвращает URL по его хешу.
	// Пример:
	//   url, err := store.GetByHash(ctx, "abc123")
	GetByHash(ctx context.Context, hash string) (URL, error)

	// Close освобождает ресурсы хранилища.
	// Пример:
	//   defer store.Close()
	Close()

	// Ping проверяет доступность хранилища.
	// Пример:
	//   if err := store.Ping(); err != nil {
	//       log.Println("Storage unavailable")
	//   }
	Ping() error

	// BatchAdd добавляет несколько URL атомарно.
	// Пример:
	//   urls := []URL{...}
	//   err := store.BatchAdd(ctx, urls, userID)
	BatchAdd(ctx context.Context, urls []URL, userID int) error

	// BatchDelete помечает URL как удаленные для указанного пользователя.
	// Пример:
	//   err := store.BatchDelete(ctx, UserHash{...})
	BatchDelete(ctx context.Context, uh UserHash) (err error)

	// GetByUserID возвращает все URL пользователя.
	// Пример:
	//   urls, err := store.GetByUserID(ctx, 1)
	GetByUserID(ctx context.Context, userID int) ([]URL, error)

	// CreateUser создает нового пользователя и возвращает его ID.
	// Пример:
	//   userID, err := store.CreateUser(ctx)
	CreateUser(ctx context.Context) (int, error)

	Stats(ctx context.Context) (model.Stats, error)
}

// CreateStore создает соответствующую реализацию Storer на основе конфигурации.
// Приоритет выбора хранилища:
//  1. PostgreSQL (если указан DSN)
//  2. Файловое хранилище (если указан путь)
//  3. In-memory хранилище (по умолчанию)
//
// Пример:
//
//	store, err := CreateStore(config)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer store.Close()
func CreateStore(conf *config.Config) (Storer, error) {
	if conf.DB.DatabaseDsn != "" {
		return newDBStore(conf)
	}
	if conf.FileStorage.FileStoragePath != "" {
		return newFileStore(conf)
	}
	return newMemStore(), nil
}
