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
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, Error("user id not found"))
			return
		}

		listID, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			log.Error("invalid id")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, Error("invalid id"))
			return
		}

		var input bookshelf.Book
		if err := render.DecodeJSON(r.Body, &input); err != nil {
			log.Error(err.Error())
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, Error("invalid request"))
			return
		}

		id, err := h.services.Book.Create(userID, listID, input)
		if err != nil {
			log.Error(err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, Error("cannot create book"))
			return
		}

		log.Info("book has been created")
		render.Status(r, http.StatusOK)
		render.JSON(w, r, createBookResponse{
			BookID: id,
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
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, Error("user id not found"))
			return
		}

		bookID, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			log.Error("invalid id")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, Error("invalid id"))
			return
		}

		books, err := h.services.Book.GetAll(userID, bookID)
		if err != nil {
			log.Error(err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, Error("cannot get books"))
			return
		}

		render.JSON(w, r, getAllBooksResponse{
			Books: books,
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
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, Error("user id not found"))
			return
		}

		bookID, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			log.Error("invalid id")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, Error("invalid id"))
			return
		}

		book, err := h.services.Book.GetByID(userID, bookID)
		if err != nil {
			log.Error(err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, Error("cannot get book"))
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, getBookByIDResponse{
			Book: book,
		})
	}
}

func (h *Handler) updateBook(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			log.Error("user id not found")
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, Error("user id not found"))
			return
		}

		bookID, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			log.Error("invalid id")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, Error("invalid id"))
			return
		}

		var input bookshelf.UpdateBookInput
		if err := render.DecodeJSON(r.Body, &input); err != nil {
			log.Error(err.Error())
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, Error("invalid request"))
			return
		}

		err = h.services.Book.Update(userID, bookID, input)
		if err != nil {
			log.Error(err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, Error("cannot update book"))
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, OK())
	}
}

func (h *Handler) deleteBook(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			log.Error("user id not found")
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, Error("user id not found"))
		}

		bookID, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			log.Error("invalid id")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, Error("invalid id"))
			return
		}
		err = h.services.Book.Delete(userID, bookID)
		if err != nil {
			log.Error(err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, err.Error())
			return
		}
		log.Info("book has been deleted")

		render.Status(r, http.StatusOK)
		render.JSON(w, r, OK())
	}
}
