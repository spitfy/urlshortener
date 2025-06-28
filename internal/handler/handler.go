package handler

import (
	"io"
	"log"
	"mime"
	"net/http"
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
	log.Println("==start==")
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println("==Content-Type==")

	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "text/plain" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	log.Println("==Content-body==")

	if err != nil || len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println("==end==")

	w.WriteHeader(http.StatusCreated)
}
