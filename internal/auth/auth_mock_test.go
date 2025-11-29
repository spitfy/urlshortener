// auth_mocks_test.go
package auth

import (
	"net/http"
	"net/http/httptest"
)

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
