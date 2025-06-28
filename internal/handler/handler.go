package handler

import (
	"fmt"
	"io"
	"net/http"
	"sync"
)

var (
	store = make(map[string]string)
	mu    sync.RWMutex
)

func Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if id := r.PathValue("id"); len(id) == 0 || len(id) > 10 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusTemporaryRedirect)
	w.Write([]byte("https://practicum.yandex.ru/"))
}

func Post(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.Header.Get("Content-Type") != "text/plain" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil || len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println(string(body))

	w.WriteHeader(http.StatusCreated)
}
