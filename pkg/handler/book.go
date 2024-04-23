package handler

import (
	bookshelf "bookshelf-api"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
)

type createBookResponse struct {
	Response
	BookID int `json:"id"`
}

func (h *Handler) createBook(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			log.Error("user id not found")
			render.JSON(w, r, Error("user id not found"))
		}

		listID, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			log.Error("invalid id")
			render.JSON(w, r, "invalid id")
			return
		}

		var input bookshelf.Book
		if err := render.DecodeJSON(r.Body, &input); err != nil {
			log.Error(err.Error())
			render.JSON(w, r, Error(err.Error()))
			return
		}

		id, err := h.services.Book.Create(userID, listID, input)
		if err != nil {
			log.Error(err.Error())
			render.JSON(w, r, Error(err.Error()))
			return
		}

		log.Info("book has been created")
		render.JSON(w, r, createBookResponse{
			Response: OK(),
			BookID:   id,
		})
	}
}

type getAllBooksResponse struct {
	Response
	Books []bookshelf.Book `json:"books"`
}

func (h *Handler) getAllBooks(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			log.Error("user id not found")
			render.JSON(w, r, Error("user id not found"))
		}

		bookID, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			log.Error("invalid id")
			render.JSON(w, r, "invalid id")
			return
		}

		books, err := h.services.Book.GetAll(userID, bookID)
		if err != nil {
			log.Error(err.Error())
			render.JSON(w, r, Error(err.Error()))
			return
		}

		render.JSON(w, r, getAllBooksResponse{
			Response: OK(),
			Books:    books,
		})
	}
}

type getBookByIDResponse struct {
	Response
	Book bookshelf.Book `json:"book"`
}

func (h *Handler) getBookByID(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			log.Error("user id not found")
			render.JSON(w, r, Error("user id not found"))
		}

		bookID, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			log.Error("invalid id")
			render.JSON(w, r, "invalid id")
			return
		}

		book, err := h.services.Book.GetByID(userID, bookID)
		if err != nil {
			log.Error(err.Error())
			render.JSON(w, r, Error(err.Error()))
			return
		}

		render.JSON(w, r, getBookByIDResponse{
			Response: OK(),
			Book:     book,
		})
	}
}

func (h *Handler) updateBook(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			log.Error("user id not found")
			render.JSON(w, r, Error("user id not found"))
		}

		bookID, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			log.Error("invalid id")
			render.JSON(w, r, "invalid id")
			return
		}

		var input bookshelf.UpdateBookInput
		if err := render.DecodeJSON(r.Body, &input); err != nil {
			log.Error(err.Error())
			render.JSON(w, r, Error(err.Error()))
			return
		}

		err = h.services.Book.Update(userID, bookID, input)
		if err != nil {
			log.Error(err.Error())
			render.JSON(w, r, Error(err.Error()))
			return
		}

		render.JSON(w, r, OK())
	}
}

func (h *Handler) deleteBook(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			log.Error("user id not found")
			render.JSON(w, r, Error("user id not found"))
		}

		bookID, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			log.Error("invalid id")
			render.JSON(w, r, "invalid id")
			return
		}
		err = h.services.Book.Delete(userID, bookID)
		if err != nil {
			log.Error(err.Error())
			render.JSON(w, r, err.Error())
			return
		}
		log.Info("book has been deleted")
		render.JSON(w, r, OK())
	}
}
