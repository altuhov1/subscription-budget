package handlers

import (
	"budget/internal/services"
	"net/http"
)

type Handler struct {
	SubManager services.SubManager
}

func NewHandler(SubManager services.SubManager) (*Handler, error) {

	return &Handler{
		SubManager: SubManager,
	}, nil
}

func (h *Handler) CreateSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Метод не разрешен. Используйте POST", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) GetSubWithParamHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "Метод не разрешен. Используйте GET", http.StatusMethodNotAllowed)
	}
}
func (h *Handler) AllSubsByUserIDHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.Header().Set("Allow", "DELETE")
		http.Error(w, "Метод не разрешен. Используйте DELETE", http.StatusMethodNotAllowed)
	}
}
func (h *Handler) GetSubsByIDHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "Метод не разрешен. Используйте GET", http.StatusMethodNotAllowed)
	}
}
func (h *Handler) UpdateSubByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.Header().Set("Allow", "PUT")
		http.Error(w, "Метод не разрешен. Используйте PUT", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) DeleteSubnByIDHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.Header().Set("Allow", "DELETE")
		http.Error(w, "Метод не разрешен. Используйте DELETE", http.StatusMethodNotAllowed)
	}
}
