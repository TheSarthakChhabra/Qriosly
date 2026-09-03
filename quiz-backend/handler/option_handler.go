package handler

import (
	"encoding/json"
	"net/http"
	"quiz-backend/service"
	"strings"
)

type OptionHandler struct {
	service *service.OptionService
}

func NewOptionHandler(s *service.OptionService) *OptionHandler {
	return &OptionHandler{service: s}
}

func (h *OptionHandler) CreateOption(w http.ResponseWriter, r *http.Request) {
	questionID := r.PathValue("id")
	var input struct {
		Text      string `json:"text"`
		IsCorrect bool   `json:"is_correct"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(input.Text) == "" {
		http.Error(w, "text is required", http.StatusBadRequest)
		return
	}

	option, err := h.service.CreateOption(r.Context(), questionID, input.Text, input.IsCorrect)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(option)
}

func (h *OptionHandler) GetOptionsByQuestionID(w http.ResponseWriter, r *http.Request) {
	questionID := r.PathValue("id")
	options, err := h.service.GetOptionsByQuestionID(r.Context(), questionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "appliaction/json")
	json.NewEncoder(w).Encode(options)
}
