package handler

import (
	bookshelf "bookshelf-api"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
)

type createListResponse struct {
	Response
	ListID int `json:"list_id"`
}

func (h *Handler) createList(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			log.Error("user id not found")
			render.JSON(w, r, Error("user id not found"))
		}

		var input bookshelf.List
		if err := render.DecodeJSON(r.Body, &input); err != nil {
			log.Error(err.Error())
			render.JSON(w, r, Error(err.Error()))
			return
		}

		id, err := h.services.List.Create(userID, input)
		if err != nil {
			log.Error(err.Error())
			render.JSON(w, r, Error(err.Error()))
			return
		}

		log.Info("list has been created")
		render.JSON(w, r, createListResponse{
			Response: OK(),
			ListID:   id,
		})
	}
}

type getAllListsResponse struct {
	Response
	Data []bookshelf.List `json:"data"`
}

func (h *Handler) getAllLists(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			log.Error("user id not found")
			render.JSON(w, r, Error("user id not found"))
		}

		lists, err := h.services.List.GetAll(userID)
		if err != nil {
			log.Error(err.Error())
			render.JSON(w, r, err.Error())
		}

		log.Info("lists have been received successfully")

		render.JSON(w, r, getAllListsResponse{
			Response: OK(),
			Data:     lists,
		})
	}
}

type getListByIDResponse struct {
	Response
	Data bookshelf.List `json:"data"`
}

func (h *Handler) getListByID(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			log.Error("user id not found")
			render.JSON(w, r, Error("user id not found"))
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			log.Error("invalid id")
			render.JSON(w, r, "invalid id")
			return
		}

		list, err := h.services.List.GetByID(userID, id)
		if err != nil {
			log.Error(err.Error())
			render.JSON(w, r, err.Error())
			return
		}

		log.Info("list have been received successfully")

		render.JSON(w, r, getListByIDResponse{
			Response: OK(),
			Data:     list,
		})
	}
}

func (h *Handler) updateList(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			log.Error("user id not found")
			render.JSON(w, r, Error("user id not found"))
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			log.Error("invalid id")
			render.JSON(w, r, "invalid id")
			return
		}

		var input bookshelf.UpdateListInput
		if err := render.DecodeJSON(r.Body, &input); err != nil {
			log.Error(err.Error())
			render.JSON(w, r, Error(err.Error()))
			return
		}

		err = h.services.List.Update(userID, id, input)
		if err != nil {
			log.Error(err.Error())
			render.JSON(w, r, Error(err.Error()))
			return
		}
		log.Info("list has been updated")
		render.JSON(w, r, OK())
	}
}

func (h *Handler) deleteList(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			log.Error("user id not found")
			render.JSON(w, r, Error("user id not found"))
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			log.Error("invalid id")
			render.JSON(w, r, "invalid id")
			return
		}
		err = h.services.List.Delete(userID, id)
		if err != nil {
			log.Error(err.Error())
			render.JSON(w, r, Error(err.Error()))
			return
		}
		log.Info("list has been deleted")
		render.JSON(w, r, OK())
	}
}
