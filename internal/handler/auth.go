package handler

import (
	"errors"
	"github.com/spitfy/urlshortener/internal/auth"
	"golang.org/x/net/context"
	"net/http"
)

func (h *Handler) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := h.auth.GetTokenFromCookie(r)
		if errors.Is(err, auth.ErrUnAuth) {
			userID, err := h.service.CreateUser(r.Context())
			if err != nil {
				http.Error(w, "error create user", http.StatusInternalServerError)
				return
			}
			tokenString, err := h.auth.BuildJWT(userID)
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
			http.SetCookie(w, &http.Cookie{
				Name:     "ID",
				Value:    tokenString,
				Path:     "/",
				HttpOnly: true,
			})
			token = tokenString
		} else if err != nil {
			http.Error(w, "invalid cookie", http.StatusUnauthorized)
			return
		}

		userID, err := h.auth.ParseUserID(token)
		if err != nil {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
