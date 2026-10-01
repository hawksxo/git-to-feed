package approval

import (
	"encoding/json"
	"errors"
	"net/http"
)

type RejectRequest struct {
	Reason string `json:"reason"`
}

type EditRequest struct {
	Content string `json:"content"`
}

type Handler struct {
	useCase *ApprovalUseCase
}

func NewHandler(useCase *ApprovalUseCase) *Handler {
	return &Handler{useCase: useCase}
}

func (h *Handler) HandleListPending(w http.ResponseWriter, r *http.Request) {
	posts, err := h.useCase.ListPending(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(posts)
}

func (h *Handler) HandleApprove(w http.ResponseWriter, r *http.Request) {
	uuid := r.PathValue("uuid")
	if uuid == "" {
		http.Error(w, "missing uuid parameter", http.StatusBadRequest)
		return
	}
	
	post, err := h.useCase.Approve(r.Context(), uuid)
	if err != nil {
		if errors.Is(err, ErrPostNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(post)
}

func (h *Handler) HandleReject(w http.ResponseWriter, r *http.Request) {
	uuid := r.PathValue("uuid")
	if uuid == "" {
		http.Error(w, "missing uuid parameter", http.StatusBadRequest)
		return
	}

	var req RejectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	post, err := h.useCase.Reject(r.Context(), uuid, req.Reason)
	if err != nil {
		if errors.Is(err, ErrPostNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(post)
}

func (h *Handler) HandleEditAndApprove(w http.ResponseWriter, r *http.Request) {
	uuid := r.PathValue("uuid")
	if uuid == "" {
		http.Error(w, "missing uuid parameter", http.StatusBadRequest)
		return
	}

	var req EditRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	post, err := h.useCase.EditAndApprove(r.Context(), uuid, req.Content)
	if err != nil {
		if errors.Is(err, ErrPostNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(post)
}
