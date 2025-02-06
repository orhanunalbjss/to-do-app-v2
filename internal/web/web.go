package web

import (
	"encoding/json"
	"net/http"
	"to-do-app-v2/internal/store"
)

type ItemError struct {
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
}

type Web struct {
	store *store.Store
}

func NewWeb(store *store.Store) *Web {
	return &Web{
		store: store,
	}
}

func (web *Web) HandleHTTPPost(w http.ResponseWriter, r *http.Request) {
	var item store.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		web.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	newItem, err := web.store.Create(item)
	if err != nil {
		web.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(newItem); err != nil {
		web.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (web *Web) HandleHTTPGet(w http.ResponseWriter, _ *http.Request) {
	items, err := web.store.ReadAll()
	if err != nil {
		web.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		web.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (web *Web) HandleHTTPGetWithId(w http.ResponseWriter, r *http.Request) {
	id := store.ItemId(r.PathValue("id"))
	item, err := web.store.Read(id)
	if err != nil {
		web.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(item)
	if err != nil {
		web.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (web *Web) HandleHTTPPut(w http.ResponseWriter, r *http.Request) {
	id := store.ItemId(r.PathValue("id"))

	var newItem store.Item
	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		web.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	item, err := web.store.Update(id, newItem)
	if err != nil {
		web.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(item); err != nil {
		web.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (web *Web) HandleHTTPDelete(w http.ResponseWriter, r *http.Request) {
	id := store.ItemId(r.PathValue("id"))

	if err := web.store.Delete(id); err != nil {
		web.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (web *Web) errorResponse(w http.ResponseWriter, statusCode int, errorString string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	encodingError := json.NewEncoder(w).Encode(ItemError{
		StatusCode: statusCode,
		Error:      errorString,
	})
	if encodingError != nil {
		http.Error(w, encodingError.Error(), http.StatusInternalServerError)
	}
}
