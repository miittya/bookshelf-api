package handler

import (
	"bookshelf-api/pkg/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
)

type Handler struct {
	services *service.Service
}

func New(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes(log *slog.Logger) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/auth", func(r chi.Router) {
		r.Post("/sign-up", h.SignUp(log))
		r.Post("/sign-in", h.SignIn(log))
	})

	router.Route("/api", func(r chi.Router) {
		r.Use(h.userIdentity(log))
		r.Route("/lists", func(r chi.Router) {
			r.Post("/", h.createList(log))
			r.Get("/", h.getAllLists(log))
			r.Get("/{id}", h.getListByID(log))
			r.Put("/{id}", h.updateList(log))
			r.Delete("/{id}", h.deleteList(log))

			r.Route("/{id}/books", func(r chi.Router) {
				r.Post("/", h.createBook(log))
				r.Get("/", h.getAllBooks(log))
			})
		})
		r.Route("/books", func(r chi.Router) {
			r.Get("/{id}", h.getBookByID(log))
			r.Put("/{id}", h.updateBook(log))
			r.Delete("/{id}", h.deleteBook(log))
		})

	})
	return router
}
