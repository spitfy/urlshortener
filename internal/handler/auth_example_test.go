package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/spitfy/urlshortener/internal/audit"
	"github.com/spitfy/urlshortener/internal/auth"
	"github.com/spitfy/urlshortener/internal/model"
	"github.com/spitfy/urlshortener/internal/repository"
	"io"
	"net/http"
	"net/http/httptest"
)

// Мок auth менеджера
type mockAuth struct{}

func (m *mockAuth) GetTokenFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("ID")
	if err != nil {
		return "", auth.ErrUnAuth
	}
	if cookie.Value == "" {
		return "", auth.ErrUnAuth
	}
	return cookie.Value, nil
}

func (m *mockAuth) ParseUserID(token string) (int, error) {
	if token == "valid-token" {
		return 42, nil
	}
	return 0, errors.New("invalid token")
}

func (m *mockAuth) CreateToken(w http.ResponseWriter, _ int) (string, error) {
	token := "valid-token"
	http.SetCookie(w, &http.Cookie{Name: "ID", Value: token, Path: "/", HttpOnly: true})
	return token, nil
}

// Мок сервиса создания пользователя
type mockService struct{}

func (m *mockService) Add(_ context.Context, _ string, _ int) (string, error) {
	return "", nil
}

func (m *mockService) BatchAdd(_ context.Context, _ []model.BatchCreateRequest, _ int) ([]model.BatchCreateResponse, error) {
	return make([]model.BatchCreateResponse, 0), nil
}

func (m *mockService) GetByHash(_ context.Context, _ string) (repository.URL, error) {
	return repository.URL{}, nil
}

func (m *mockService) Ping() error {
	return nil
}

func (m *mockService) GetByUserID(_ context.Context, _ int) ([]model.LinkPair, error) {
	return make([]model.LinkPair, 0), nil
}

func (m *mockService) DeleteEnqueue(_ context.Context, _ []string, _ int) {
}

func (m *mockService) AddObserver(_ audit.Observer) {
}

func (m *mockService) NotifyObservers(_ context.Context, _ audit.Event) {
}

func (m *mockService) CreateUser(_ context.Context) (int, error) {
	return 100, nil
}

func (m *mockService) Stats(_ context.Context) (model.Stats, error) {
	return model.Stats{URLs: 1, Users: 1}, nil
}

// Example-функция для authMiddleware
func ExampleHandler_authMiddleware() {
	h := &Handler{
		auth:    &mockAuth{},
		service: &mockService{},
	}

	protectedHandler := h.authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int)
		fmt.Fprintf(w, "Authenticated user ID: %d", userID)
	}))

	// Тестируем запрос без cookie — создается пользователь и токен
	req := httptest.NewRequest("GET", "http://example.com", nil)
	w := httptest.NewRecorder()
	protectedHandler.ServeHTTP(w, req)
	res := w.Result()

	body1, _ := io.ReadAll(res.Body)
	res.Body.Close()

	fmt.Println("Status Code:", res.StatusCode)
	fmt.Println("Body:", string(body1))

	// Тестируем запрос с валидным токеном в cookie
	req2 := httptest.NewRequest("GET", "http://example.com", nil)
	req2.AddCookie(&http.Cookie{Name: "ID", Value: "valid-token"})
	w2 := httptest.NewRecorder()
	protectedHandler.ServeHTTP(w2, req2)
	res2 := w2.Result()

	body2, _ := io.ReadAll(res2.Body)
	res2.Body.Close()

	fmt.Println("Status Code:", res2.StatusCode)
	fmt.Println("Body:", string(body2))

	// Output:
	// Status Code: 200
	// Body: Authenticated user ID: 100
	// Status Code: 200
	// Body: Authenticated user ID: 42
}
