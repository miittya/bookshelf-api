package handler

import (
	"context"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strings"
)

const (
	authHeader = "Authorization"
)

func (h *Handler) userIdentity(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get(authHeader)
			if header == "" {
				log.Error("empty auth header")
				render.JSON(w, r, Error("empty auth header"))
				return
			}

			headerParts := strings.Split(header, " ")
			if len(headerParts) != 2 {
				log.Error("invalid auth header")
				render.JSON(w, r, Error("invalid auth header"))
				return
			}

			id, err := h.services.Authorization.ParseToken(headerParts[1])
			if err != nil {
				log.Error(err.Error())
				render.JSON(w, r, Error(err.Error()))
				return
			}

			ctx := context.WithValue(r.Context(), "userID", id)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
