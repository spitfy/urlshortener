package repository

import (
	"github.com/spitfy/urlshortener/internal/config"
	repoConf "github.com/spitfy/urlshortener/internal/repository/config"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
)

var cfg = config.Config{
	FileStorage: repoConf.Config{
		FileStorageName: config.DefaultFileName,
		FileStoragePath: config.DefaultFileStorage,
	},
}

func TestNewStore(t *testing.T) {
	fullPath := filepath.Join(cfg.FileStorage.FileStoragePath, cfg.FileStorage.FileStorageName)
	f, _ := os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE, 0666)
	var store = &Store{
		mux:  &sync.RWMutex{},
		s:    map[string]link{"ASDQWE23": {"https://github.com/", "1"}},
		file: f,
	}

	tests := []struct {
		name string
		want *Store
	}{
		{"success", store},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStore(&cfg); !reflect.DeepEqual(got.s, tt.want.s) {
				t.Errorf("NewStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStore_Add(t *testing.T) {
	fullPath := filepath.Join(cfg.FileStorage.FileStoragePath, cfg.FileStorage.FileStorageName)
	f, _ := os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE, 0666)
	var store = &Store{
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
}
