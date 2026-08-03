package materials

import (
	"encoding/json"
	"errors"
	"log"
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
	w.Header().Set("Content-Type", "application/json")

	var body addRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Println(err)
		http.Error(w, "decode body error", http.StatusBadRequest)
		return
	}

	id, err := h.service.Add(r.Context(), body.CategoryID, body.Price, body.ImageURL)
	if err != nil {
		if errors.Is(err, myerrors.ErrBadRequest) {
			log.Println(err)
			http.Error(w, "bad request", http.StatusBadRequest)
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
		http.Error(w, "encode body error", http.StatusInternalServerError)
		return
	}
	log.Println("Запрос на POST /materials был успешно обработан")
}

func (h *Handler) GetByCategoryID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	materials, err := h.service.GetByCategoryID(r.Context(), id)
	if err != nil {
		if errors.Is(err, myerrors.ErrBadRequest) {
			log.Println(err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		log.Println(err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(materials)
	if err != nil {
		log.Println(err)
		http.Error(w, "encode body error", http.StatusInternalServerError)
		return
	}

	log.Println("Запрос на GET /materials/{id} был успешно обработан")
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, myerrors.ErrNotFound) {
			log.Println(err)
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, myerrors.ErrBadRequest) {
			log.Println(err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		log.Println(err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	log.Println("Запрос на DELETE успешно обработан")
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, adminAuth func(http.HandlerFunc) http.HandlerFunc) {
	mux.HandleFunc("POST /materials", adminAuth(h.Add))
	mux.HandleFunc("GET /materials/{id}", h.GetByCategoryID)
	mux.HandleFunc("DELETE /materials/{id}", adminAuth(h.Delete))
}
