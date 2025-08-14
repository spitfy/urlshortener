package handler

import (
	"errors"
	"github.com/spitfy/urlshortener/internal/auth"
	"golang.org/x/net/context"
	"net/http"
)

func (h *Handler) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userID int
		token, err := h.auth.GetTokenFromCookie(r)
		if errors.Is(err, auth.ErrUnAuth) {
			if cookies := r.Cookies(); len(cookies) > 0 {
				if _, err := h.initCookie(w, r); err != nil {
					http.Error(w, "error create user", http.StatusInternalServerError)
					return
				}
				http.Error(w, "", http.StatusNoContent)
				return
			}
			token, err = h.initCookie(w, r)
			if err != nil {
				http.Error(w, "error create user", http.StatusInternalServerError)
				return
			}
		} else if err != nil {
			http.Error(w, "invalid cookie", http.StatusUnauthorized)
			return
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

func (h *Handler) initCookie(w http.ResponseWriter, r *http.Request) (string, error) {
	userID, err := h.service.CreateUser(r.Context())
	if err != nil {
		return "", err
	}
	return h.auth.CreateToken(w, userID)
}
