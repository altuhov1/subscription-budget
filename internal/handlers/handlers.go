package handlers

import (
	"budget/internal/models"
	"budget/internal/services"
	"encoding/json"
	"log/slog"
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
		http.Error(w, "The method is not allowed. Use the POST", http.StatusMethodNotAllowed)
		return
	}

	var req models.Subscription
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("Invalid JSON in link request", "error", err)
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	NumOfSub, err := h.SubManager.CreateSubscription(ctx, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(struct {
		NumOfSub int `json:"numberOfSub"`
	}{
		NumOfSub: NumOfSub,
	})
}

func (h *Handler) GetSubWithParamHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "Метод не разрешен. Используйте GET", http.StatusMethodNotAllowed)
		return
	}

	var req models.TotalCostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("Invalid JSON in total cost request", "error", err)
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	resp, err := h.SubManager.GetSubscriptionsWithParam(ctx, &req)
	if err != nil {
		slog.Error("Failed to calculate total cost", "error", err)
		http.Error(w, `{"error":"calculation failed"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
func (h *Handler) AllSubsByUserIDHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "The method is not allowed. Use the GET", http.StatusMethodNotAllowed)
		return
	}

	type Request struct {
		UserID string `json:"user_id"`
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("Invalid JSON in list request", "error", err)
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}

	if req.UserID == "" {
		http.Error(w, `{"error":"user_id is required"}`, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	subs, err := h.SubManager.ListSubscriptionsByUserID(ctx, req.UserID)
	if err != nil {
		slog.Error("Failed to list subscriptions", "user_id", req.UserID, "error", err)
		http.Error(w, `{"error":"failed to retrieve subscriptions"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subs)
}
func (h *Handler) GetSubsByIDHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "The method is not allowed. Use the GET", http.StatusMethodNotAllowed)
		return
	}

	type Request struct {
		ID int `json:"id"`
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("Invalid JSON in get by ID request", "error", err)
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}

	if req.ID <= 0 {
		http.Error(w, `{"error":"invalid ID"}`, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	sub, err := h.SubManager.GetSubscriptionByID(ctx, req.ID)
	if err != nil {
		slog.Error("Failed to get subscription", "id", req.ID, "error", err)
		http.Error(w, `{"error":"subscription not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sub)
}

func (h *Handler) UpdateSubByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.Header().Set("Allow", "PUT")
		http.Error(w, "The method is not allowed. Use the PUT", http.StatusMethodNotAllowed)
		return
	}

	var req models.SubscriptionWithId
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("Invalid JSON in update request", "error", err)
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}

	if req.ID <= 0 {
		http.Error(w, `{"error":"invalid ID"}`, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.SubManager.UpdateSubscriptionByID(ctx, &req); err != nil {
		slog.Error("Failed to update subscription", "id", req.ID, "error", err)
		http.Error(w, `{"error":"update failed"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteSubnByIDHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.Header().Set("Allow", "DELETE")
		http.Error(w, "The method is not allowed. Use the DELETE", http.StatusMethodNotAllowed)
		return
	}

	type Request struct {
		ID int `json:"id"`
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("Invalid JSON in delete request", "error", err)
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}

	if req.ID <= 0 {
		http.Error(w, `{"error":"invalid ID"}`, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.SubManager.DeleteSubscriptionByID(ctx, req.ID); err != nil {
		slog.Error("Failed to delete subscription", "id", req.ID, "error", err)
		http.Error(w, `{"error":"deletion failed"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
