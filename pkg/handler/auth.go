package handler

import (
	bookshelf "bookshelf-api"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
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
		if err != nil {
			log.Error("invalid request")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, Error("invalid request"))
			return
		}

		if err := validator.New().Struct(input); err != nil {
			log.Error("invalid request")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, Error("invalid request"))
			return
		}
		id, err := h.services.Authorization.CreateUser(input)
		if err != nil {
			log.Error(err.Error())

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, Error("cannot create user"))
			return
		}

		log.Info("user has been created")
		render.JSON(w, r, signUpResponse{
			ID: id,
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
		if err != nil {
			log.Error("invalid request", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, Error("invalid request"))

			return
		}
		token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
		if err != nil {
			log.Error(err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, Error("cant generate token"))
			return
		}

		log.Info("token has been generated")
		render.JSON(w, r, signInResponse{
			Token: token,
		})
	}
}
