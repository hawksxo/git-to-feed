package pipeline

import (
	"encoding/json"
	"net/http"
)

type AdminHandler struct {
	repo FewShotRepository
}

func NewAdminHandler(repo FewShotRepository) *AdminHandler {
	return &AdminHandler{repo: repo}
}

func (h *AdminHandler) HandleListExamples(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	examples, err := h.repo.FindAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(examples)
}

func (h *AdminHandler) HandleCreateExample(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var example FewShotExampleEntity
	if err := json.NewDecoder(r.Body).Decode(&example); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if err := h.repo.Save(r.Context(), &example); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(example)
}
