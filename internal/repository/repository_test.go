package repository

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"sync"
	"testing"
)

func TestNewStore(t *testing.T) {
	var store = &Store{
		mux: &sync.Mutex{},
		s:   make(map[string]string),
	}

	tests := []struct {
		name string
		want *Store
	}{
		{"success", store},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStore(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStore_Add(t *testing.T) {
	store := &Store{
		mux: &sync.Mutex{},
		s:   make(map[string]string),
	}
	tests := []struct {
		name string
		link URL
		want string
	}{
		{"success", URL{"https://github.com/", "ASDQWE23"}, "https://github.com/"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store.Add(tt.link)
			assert.Equal(t, tt.want, store.s[tt.link.Hash])
		})
	}
}
