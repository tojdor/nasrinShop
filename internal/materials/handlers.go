package materials

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	myerrors "github.com/tojdor/nasrinShop/internal/myErrors"
)

type Handler struct {
	service *Service
}

type addRequest struct {
	CategoryID int    `json:"category_id"`
	Price      int    `json:"price"`
	ImageURL   string `json:"image_url"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {
	var body addRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "decode body error", http.StatusBadRequest)
		return
	}

	id, err := h.service.Add(r.Context(), body.CategoryID, body.Price, body.ImageURL)
	if err != nil {
		if errors.Is(err, myerrors.ErrBadRequest) {
			http.Error(w, "bad request", http.StatusBadRequest)
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
		http.Error(w, "encode body error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetByCategorieID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	materials, err := h.service.GetByCategorieID(r.Context(), id)
	if err != nil {
		if errors.Is(err, myerrors.ErrBadRequest) {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(materials)
	if err != nil {
		http.Error(w, "encode body error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, myerrors.ErrNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, myerrors.ErrBadRequest) {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /material", h.Add)
	mux.HandleFunc("GET /materials/{id}", h.GetByCategorieID)
	mux.HandleFunc("DELETE /material/{id}", h.Delete)
}
