package handler

import (
	bookshelf "bookshelf-api"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
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
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, Error("user id not found"))
			return
		}

		var input bookshelf.List
		if err := render.DecodeJSON(r.Body, &input); err != nil {
			log.Error(err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, Error("invalid request"))
			return
		}
		if err := validator.New().Struct(input); err != nil {
			log.Error(err.Error())
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, Error("invalid request"))
			return
		}

		id, err := h.services.List.Create(userID, input)
		if err != nil {
			log.Error(err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, Error("cannot create list"))
			return
		}

		log.Info("list has been created")
		render.Status(r, http.StatusOK)
		render.JSON(w, r, createListResponse{
			ListID: id,
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
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, Error("user id not found"))
			return
		}

		lists, err := h.services.List.GetAll(userID)
		if err != nil {
			log.Error(err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, Error("cannot get all lists"))
			return
		}

		log.Info("lists have been received successfully")
		render.Status(r, http.StatusOK)
		render.JSON(w, r, getAllListsResponse{
			Data: lists,
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
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, Error("user id not found"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			log.Error("invalid id")
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, Error("invalid id"))
			return
		}

		list, err := h.services.List.GetByID(userID, id)
		if err != nil {
			log.Error(err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, Error("cannot get list by id"))
			return
		}

		log.Info("list have been received successfully")

		render.Status(r, http.StatusOK)
		render.JSON(w, r, getListByIDResponse{
			Data: list,
		})
	}
}

func (h *Handler) updateList(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			log.Error("user id not found")
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, Error("user id not found"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			log.Error("invalid id")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, Error("invalid id"))
			return
		}

		var input bookshelf.UpdateListInput
		if err := render.DecodeJSON(r.Body, &input); err != nil {
			log.Error(err.Error())
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, Error("invalid request"))
			return
		}

		err = h.services.List.Update(userID, id, input)
		if err != nil {
			log.Error(err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, Error("cannot update list"))
			return
		}
		log.Info("list has been updated")
		render.Status(r, http.StatusOK)
		render.JSON(w, r, OK())
	}
}

func (h *Handler) deleteList(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			log.Error("user id not found")
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, Error("user id not found"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			log.Error("invalid id")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, Error("invalid id"))
			return
		}
		err = h.services.List.Delete(userID, id)
		if err != nil {
			log.Error(err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, Error(err.Error()))
			return
		}
		log.Info("list has been deleted")
		render.Status(r, http.StatusOK)
		render.JSON(w, r, OK())
	}
}
