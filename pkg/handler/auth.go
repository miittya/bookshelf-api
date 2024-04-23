package handler

import (
	bookshelf "bookshelf-api"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"net/http"
)

type signUpResponse struct {
	Response
	ID int `json:"id"`
}

type signInResponse struct {
	Response
	Token string `json:"token"`
}

func (h *Handler) SignUp(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var input bookshelf.User

		err := render.DecodeJSON(r.Body, &input)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})

			render.JSON(w, r, Error("failed to decode request"))

			return
		}
		id, err := h.services.Authorization.CreateUser(input)
		if err != nil {
			log.Error(err.Error())
			render.JSON(w, r, Error(err.Error()))
			return
		}

		log.Info("user has been created")
		render.JSON(w, r, signUpResponse{
			Response: OK(),
			ID:       id,
		})
	}
}

func (h *Handler) SignIn(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var input bookshelf.User

		err := render.DecodeJSON(r.Body, &input)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})

			render.JSON(w, r, Error("failed to decode request"))

			return
		}
		token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
		if err != nil {
			log.Error(err.Error())
			render.JSON(w, r, Error(err.Error()))
			return
		}

		log.Info("user has been created")
		render.JSON(w, r, signInResponse{
			Response: OK(),
			Token:    token,
		})
	}
}
