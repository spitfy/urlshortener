package handler

import (
	"encoding/json"
	models "github.com/spitfy/urlshortener/internal/model"
	"io"
	"log"
	"mime"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/spitfy/urlshortener/internal/service"
)

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hash := chi.URLParam(r, "hash")
	if len(hash) == 0 || len(hash) > service.CharCnt {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	link, err := h.service.Get(hash)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Location", link)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || !allowedContent[mediaType] {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer func() {
		_ = r.Body.Close()
	}()

	if err != nil || len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortURL, err := h.service.Add(string(body))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error saving url: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(shortURL))
}

func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req models.Request
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	link, err := h.service.Add(req.URL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	res := models.Response{Result: link}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Batch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req []models.BatchCreateRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	batchResponse, err := h.service.BatchAdd(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(batchResponse); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Ping(w http.ResponseWriter, _ *http.Request) {
	if err := h.service.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
