package repository

import (
	"github.com/spitfy/urlshortener/internal/config"
	repoConf "github.com/spitfy/urlshortener/internal/repository/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"reflect"
	"sync"
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
		mux:  &sync.RWMutex{},
		s:    make(map[string]link),
		file: f,
	}
	tests := []struct {
		name string
		link URL
		want link
	}{
		{"success", URL{"https://github.com/", "ASDQWE23"}, link{"https://github.com/", "1"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = store.Add(tt.link)
			assert.Equal(t, tt.want, store.s[tt.link.Hash])
		})
	}
	if err := os.Remove(cfg.FileStorage.FileStoragePath); err != nil {
		assert.NoError(t, err)
	}
}

func TestNewStore(t *testing.T) {
	f, _ := os.OpenFile(cfg.FileStorage.FileStoragePath, os.O_RDWR|os.O_CREATE, 0666)
	var store = &FileStore{
		mux:  &sync.RWMutex{},
		s:    make(map[string]link, 0),
		file: f,
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
