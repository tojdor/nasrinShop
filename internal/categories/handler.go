package categories

import (
	"encoding/json"
	"errors"
	"log"
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
	w.Header().Set("Content-Type", "application/json")

	name := r.PathValue("name")
	id, err := h.service.Add(r.Context(), name)
	if err != nil {
		if errors.Is(err, myerrors.ErrBadRequest) {
			log.Println(err)
			http.Error(w, "error bad request", http.StatusBadRequest)
			return
		}
		if errors.Is(err, myerrors.ErrConflict) {
			log.Println(err)
			http.Error(w, "conflict", http.StatusConflict)
			return
		}
		log.Println(err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(struct {
		ID int `json:"id"`
	}{
		ID: id,
	})
	if err != nil {
		log.Println(err)
		http.Error(w, "error json encode", http.StatusInternalServerError)
		return
	}
	log.Println("Запрос на POST /categories был успешно обработан")
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	categories, err := h.service.GetAll(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(categories)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	log.Println("Запрос на GET /categories был успешно обработан")
}

func (h *Handler) GetIDByName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	name := r.PathValue("name")
	id, err := h.service.GetIDByName(r.Context(), name)
	if err != nil {
		if errors.Is(err, myerrors.ErrBadRequest) {
			log.Println(err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		if errors.Is(err, myerrors.ErrNotFound) {
			log.Println(err)
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		log.Println(err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(struct {
		ID int `json:"id"`
	}{
		ID: id,
	})
	if err != nil {
		log.Println(err)
		http.Error(w, "encode error", http.StatusInternalServerError)
		return
	}

	log.Println("Запрос на GET /categories/{name} был успешно обработан")
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	name := r.PathValue("name")
	err := h.service.Delete(r.Context(), name)
	if err != nil {
		if errors.Is(err, myerrors.ErrBadRequest) {
			log.Println(err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		if errors.Is(err, myerrors.ErrNotFound) {
			log.Println(err)
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		log.Println(err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	log.Println("Запрос на DELETE /categories/{name} был успешно обработан")
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, adminAuth func(http.HandlerFunc) http.HandlerFunc) {
	mux.HandleFunc("GET /categories", h.GetAll)
	mux.HandleFunc("POST /categories/{name}", h.Add)
	mux.HandleFunc("GET /categories/id/{name}", h.GetIDByName)
	mux.HandleFunc("DELETE /categories/{name}", h.Delete)
}
