package categories

import (
	"encoding/json"
	"errors"
	"net/http"

	myerrors "github.com/tojdor/nasrinShop/internal/myErrors"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	id, err := h.service.Add(r.Context(), name)
	if err != nil {
		if errors.Is(err, myerrors.ErrBadRequest) {
			http.Error(w, "error bad request", http.StatusBadRequest)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(struct {
		ID int `json:"id"`
	}{
		ID: id,
	})
	if err != nil {
		http.Error(w, "error json encode", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetAll(r.Context())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(categories)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

}

func (h *Handler) GetIDByName(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	id, err := h.service.GetIDByName(r.Context(), name)
	if err != nil {
		if errors.Is(err, myerrors.ErrBadRequest) {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		if errors.Is(err, myerrors.ErrNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(struct {
		ID int `json:"id"`
	}{
		ID: id,
	})
	if err!=nil{
		http.Error(w, "encode error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	err := h.service.Delete(r.Context(), name)
	if err != nil {
		if errors.Is(err, myerrors.ErrBadRequest) {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		if errors.Is(err, myerrors.ErrNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /category", h.GetAll)
	mux.HandleFunc("POST /category", h.Add)
	mux.HandleFunc("GET /category/id/{name}", h.GetIDByName)
	mux.HandleFunc("DELETE /category/{name}", h.Delete)
}
