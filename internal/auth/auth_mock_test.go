// auth_mocks_test.go
package auth

import (
	"net/http"
	"net/http/httptest"
)

// MockResponseWriter для тестирования записи cookie
type MockResponseWriter struct {
	HeaderMap http.Header
	Body      []byte
	Status    int
}

func NewMockResponseWriter() *MockResponseWriter {
	return &MockResponseWriter{
		HeaderMap: make(http.Header),
	}
}

func (m *MockResponseWriter) Header() http.Header {
	return m.HeaderMap
}

func (m *MockResponseWriter) Write(b []byte) (int, error) {
	m.Body = append(m.Body, b...)
	return len(b), nil
}

func (m *MockResponseWriter) WriteHeader(statusCode int) {
	m.Status = statusCode
}

// MockRequest с cookie для тестирования
func NewMockRequestWithCookie(token string) *http.Request {
	req := httptest.NewRequest("GET", "/", nil)
	if token != "" {
		req.AddCookie(&http.Cookie{Name: "ID", Value: token})
	}
	return req
}

// MockRequest без cookie для тестирования
func NewMockRequestWithoutCookie() *http.Request {
	return httptest.NewRequest("GET", "/", nil)
}
