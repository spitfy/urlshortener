package service

import (
	"testing"

	"github.com/spitfy/urlshortener/internal/config"
	"github.com/spitfy/urlshortener/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestService_makeURL(t *testing.T) {
	s := &Service{
		store:  *repository.NewStore(),
		config: config.GetConfig(false),
	}

	tests := []struct {
		name  string
		hash  string
		want  string
		isErr bool
	}{
		{"sucsess", "ABCDEFGI", s.config.Service.ServerURL + "/ABCDEFGI", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			link, err := s.makeURL(tt.hash)

			assert.NoError(t, err, "Error making url")
			assert.Equal(t, tt.want, link, "Wrong url")
		})
	}
}
