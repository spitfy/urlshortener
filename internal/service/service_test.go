package service

import (
	"fmt"
	"testing"

	serviceConf "github.com/spitfy/urlshortener/internal/service/config"

	"github.com/spitfy/urlshortener/internal/config"
	"github.com/spitfy/urlshortener/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestService_makeURL(t *testing.T) {
	cfg := config.Config{Service: serviceConf.Config{ServerURL: config.DefaultServerURL}}
	store := repository.NewMockStore()
	s := &Service{
		store:  store,
		config: cfg,
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

func TestRandString(t *testing.T) {
	tests := []struct {
		name string
		len  int
		want int
	}{
		{"Six len", 6, 6},
		{"Eight len", 8, 8},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RandString(tt.len)
			assert.Equal(t, tt.want, len(got), "Wrong result")
		})
	}
}

func Test_isURL(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"success", args{"http://google.com"}, true},
		{"fail", args{"google.com"}, false},
		{"empty", args{" "}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, isURL(tt.args.str), "isURL(%v)", tt.args.str)
		})
	}
}

func BenchmarkRandString(b *testing.B) {
	lengths := []int{5, 10, 32, 64}
	for _, n := range lengths {
		b.Run(fmt.Sprintf("Len-%d", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = RandString(n)
			}
		})
	}
}
