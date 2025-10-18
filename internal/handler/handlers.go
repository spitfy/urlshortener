package handler

import (
	"encoding/json"
	"errors"

	"github.com/spitfy/urlshortener/internal/audit"
	"github.com/spitfy/urlshortener/internal/model"
	"github.com/spitfy/urlshortener/internal/repository"

	"io"
	"log"
	"mime"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/spitfy/urlshortener/internal/service"
)

// Get обрабатывает запрос на получение оригинального URL по хешу
// @Summary Получить оригинальный URL
// @Description Перенаправляет на оригинальный URL по сокращенному хешу
// @Tags URL
// @Param hash path string true "Хеш сокращенного URL"
// @Success 307 {string} string "Перенаправление на оригинальный URL"
// @Success 410 {string} string "URL был удален"
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 401 {string} string "Неавторизованный доступ"
// @Router /{hash} [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	hash := chi.URLParam(r, "hash")
	if len(hash) == 0 || len(hash) > service.CharCnt {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u, err := h.service.GetByHash(r.Context(), hash)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	if u.DeletedFlag {
		w.WriteHeader(http.StatusGone)
		return
	}

	h.service.NotifyObservers(r.Context(), audit.Event{
		Timestamp: time.Now(),
		Action:    audit.Follow,
		UserID:    userID,
		URL:       u.Link,
	})

	w.Header().Add("Location", u.Link)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// GetByUserID возвращает все сокращенные URL пользователя
// @Summary Получить URL пользователя
// @Description Возвращает все сокращенные URL, созданные текущим пользователем
// @Tags User
// @Produce json
// @Success 200 {array} model.URL "Список сокращенных URL"
// @Success 204 {string} string "Нет сохраненных URL"
// @Failure 401 {string} string "Неавторизованный доступ"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /api/user/urls [get]
func (h *Handler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	res, err := h.service.GetByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(res) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "encoding error", http.StatusInternalServerError)
		return
	}
}

// Post создает новый сокращенный URL из текстового тела запроса
// @Summary Сократить URL (текст)
// @Description Создает новый сокращенный URL из текстового тела запроса
// @Tags URL
// @Accept text/plain
// @Produce text/plain
// @Param url body string true "Оригинальный URL для сокращения"
// @Success 201 {string} string "Создан новый сокращенный URL"
// @Success 409 {string} string "URL уже был сокращен ранее"
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 401 {string} string "Неавторизованный доступ"
// @Router / [post]
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

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	shortURL, err := h.service.Add(r.Context(), string(body), userID)

	if err != nil {
		if errors.Is(err, repository.ErrExistsURL) {
			w.WriteHeader(http.StatusConflict)
			_, _ = w.Write([]byte(shortURL))
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error saving url: %v", err)
		return
	}

	h.service.NotifyObservers(r.Context(), audit.Event{
		Timestamp: time.Now(),
		Action:    audit.Shorten,
		UserID:    userID,
		URL:       string(body),
	})

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(shortURL))
}

// ShortenURL создает новый сокращенный URL из JSON тела запроса
// @Summary Сократить URL (JSON)
// @Description Создает новый сокращенный URL из JSON тела запроса
// @Tags URL
// @Accept json
// @Produce json
// @Param request body model.Request true "Запрос на сокращение URL"
// @Success 201 {object} model.Response "Создан новый сокращенный URL"
// @Success 409 {object} model.Response "URL уже был сокращен ранее"
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 401 {string} string "Неавторизованный доступ"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /api/shorten [post]
func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		http.Error(w, "invalid content-type", http.StatusBadRequest)
		return
	}

	var req model.Request
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	shortURL, err := h.service.Add(r.Context(), req.URL, userID)
	res := model.Response{Result: shortURL}
	w.Header().Set("Content-Type", "application/json")

	switch {
	case err == nil:
		w.WriteHeader(http.StatusCreated)
	case errors.Is(err, repository.ErrExistsURL):
		w.WriteHeader(http.StatusConflict)
	default:
		http.Error(w, "could not shorten URL", http.StatusInternalServerError)
		return
	}

	h.service.NotifyObservers(r.Context(), audit.Event{
		Timestamp: time.Now(),
		Action:    audit.Shorten,
		UserID:    userID,
		URL:       req.URL,
	})

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "encoding error", http.StatusInternalServerError)
		return
	}
}

// BatchAdd создает несколько сокращенных URL из пакетного запроса
// @Summary Пакетное сокращение URL
// @Description Создает несколько сокращенных URL из пакетного запроса
// @Tags URL
// @Accept json
// @Produce json
// @Param request body []model.BatchCreateRequest true "Список URL для сокращения"
// @Success 201 {array} model.BatchCreateResponse "Список созданных сокращенных URL"
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 401 {string} string "Неавторизованный доступ"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /api/shorten/batch [post]
func (h *Handler) BatchAdd(w http.ResponseWriter, r *http.Request) {
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

	var req []model.BatchCreateRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	batchResponse, err := h.service.BatchAdd(r.Context(), req, userID)
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

// Delete помечает URL как удаленные
// @Summary Удалить URL
// @Description Помечает указанные URL как удаленные (асинхронно)
// @Tags URL
// @Accept json
// @Produce json
// @Param request body []string true "Список хешей URL для удаления"
// @Success 202 {string} string "Запрос на удаление принят"
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 401 {string} string "Неавторизованный доступ"
// @Router /api/user/urls [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.Header().Set("Allow", http.MethodDelete)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req []string
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	h.service.DeleteEnqueue(r.Context(), req, userID)
	w.WriteHeader(http.StatusAccepted)
}

// Ping проверяет доступность сервера
// @Summary Проверить доступность сервера
// @Description Проверяет, что сервер работает и доступен
// @Tags Health
// @Success 200 {string} string "Сервер доступен"
// @Failure 500 {string} string "Сервер недоступен"
// @Router /ping [get]
func (h *Handler) Ping(w http.ResponseWriter, _ *http.Request) {
	if err := h.service.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
