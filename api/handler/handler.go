package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"to-do-app-v2/internal/store"
)

type Service interface {
	Create(item store.Item) (store.Item, error)
	ReadAll() ([]store.Item, error)
	Read(id store.ItemID) (store.Item, error)
	Update(id store.ItemID, item store.Item) (store.Item, error)
	Delete(id store.ItemID) error
}

type Error struct {
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) HandleHTTPPost(w http.ResponseWriter, r *http.Request) {
	var item store.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		slog.ErrorContext(r.Context(), err.Error())
		h.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	newItem, err := h.service.Create(item)
	if err != nil {
		slog.ErrorContext(r.Context(), err.Error())
		h.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(newItem); err != nil {
		slog.ErrorContext(r.Context(), err.Error())
		h.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *Handler) HandleHTTPGet(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.ReadAll()
	if err != nil {
		slog.ErrorContext(r.Context(), err.Error())
		h.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		slog.ErrorContext(r.Context(), err.Error())
		h.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *Handler) HandleHTTPGetWithID(w http.ResponseWriter, r *http.Request) {
	id := store.ItemID(r.PathValue("id"))
	item, err := h.service.Read(id)
	if err != nil {
		slog.ErrorContext(r.Context(), err.Error())
		h.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(item)
	if err != nil {
		slog.ErrorContext(r.Context(), err.Error())
		h.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *Handler) HandleHTTPPut(w http.ResponseWriter, r *http.Request) {
	id := store.ItemID(r.PathValue("id"))

	var newItem store.Item
	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		slog.ErrorContext(r.Context(), err.Error())
		h.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	item, err := h.service.Update(id, newItem)
	if err != nil {
		slog.ErrorContext(r.Context(), err.Error())
		h.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(item); err != nil {
		slog.ErrorContext(r.Context(), err.Error())
		h.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *Handler) HandleHTTPDelete(w http.ResponseWriter, r *http.Request) {
	id := store.ItemID(r.PathValue("id"))

	if err := h.service.Delete(id); err != nil {
		slog.ErrorContext(r.Context(), err.Error())
		h.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) errorResponse(w http.ResponseWriter, statusCode int, errorString string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	encodingError := json.NewEncoder(w).Encode(Error{
		StatusCode: statusCode,
		Error:      errorString,
	})
	if encodingError != nil {
		http.Error(w, encodingError.Error(), http.StatusInternalServerError)
	}
}
