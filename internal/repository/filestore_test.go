package repository

import (
	"context"
	"github.com/spitfy/urlshortener/internal/config"
	repoConf "github.com/spitfy/urlshortener/internal/repository/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"reflect"
	"testing"
)

var cfg = config.Config{
	FileStorage: repoConf.Config{FileStoragePath: config.DefaultFileStorageTest},
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestStore_Add(t *testing.T) {
	f, _ := os.OpenFile(cfg.FileStorage.FileStoragePath, os.O_RDWR|os.O_CREATE, 0666)
	var store = &FileStore{
		file:     f,
		MemStore: newMemStore(),
	}
	tests := []struct {
		name string
		link URL
		want map[string]string
	}{
		{"success", URL{"https://github.com/", "ASDQWE23"}, map[string]string{"ASDQWE23": "https://github.com/"}},
	}
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _ = store.Add(ctx, tt.link, -1)
			assert.Equal(t, tt.want[tt.link.Hash], store.s[tt.link.Hash])
		})
	}
	if err := os.Remove(cfg.FileStorage.FileStoragePath); err != nil {
		assert.NoError(t, err)
	}
}

func TestNewStore(t *testing.T) {
	f, _ := os.OpenFile(cfg.FileStorage.FileStoragePath, os.O_RDWR|os.O_CREATE, 0666)
	var store = &FileStore{
		file:     f,
		MemStore: newMemStore(),
	}

	tests := []struct {
		name string
		want *FileStore
	}{
		{"success", store},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := newFileStore(&cfg); !reflect.DeepEqual((*got).s, tt.want.s) {
				require.NoError(t, err, "error creating store")
				t.Errorf("NewFileStore() = %v, want %v", got, tt.want)
			}
		})
	}
}
