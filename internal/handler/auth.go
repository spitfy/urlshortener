package handler

import (
	"errors"
	"net/http"

	"github.com/spitfy/urlshortener/internal/auth"
	"golang.org/x/net/context"
)

func (h *Handler) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userID int
		var token string
		var err error

		token, err = h.auth.GetTokenFromCookie(r)
		if err != nil {
			if errors.Is(err, auth.ErrUnAuth) {
				userID, token, err = h.createUserAndToken(w, r)
				if err != nil {
					http.Error(w, "error create user", http.StatusInternalServerError)
					return
				}
			} else {
				http.Error(w, "invalid cookie", http.StatusUnauthorized)
				return
			}
		}

		if userID == 0 {
			userID, err = h.auth.ParseUserID(token)
			if err != nil {
				http.Error(w, "", http.StatusUnauthorized)
				return
			}
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (h *Handler) createUserAndToken(w http.ResponseWriter, r *http.Request) (int, string, error) {
	userID, err := h.service.CreateUser(r.Context())
	if err != nil {
		return 0, "", err
	}
	token, err := h.auth.CreateToken(w, userID)
	return userID, token, err
}
