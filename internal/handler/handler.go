package handler

import (
	"io"
	"log"
	"mime"
	"net/http"
)

var store []byte

func Get(w http.ResponseWriter, r *http.Request) {
	log.Println("==start==")
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println("==id==")
	if id := r.PathValue("id"); len(id) == 0 || len(id) > 10 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println("==after==", string(store))
	w.Header().Add("Location", string(store))
	w.WriteHeader(http.StatusTemporaryRedirect)
	w.Write(store)
}

func Post(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "text/plain" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil || len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	store = body

	w.WriteHeader(http.StatusCreated)
}
