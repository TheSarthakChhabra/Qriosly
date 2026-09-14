package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"quiz-backend/service"
)

type AnswerHandler struct {
	service *service.AnswerService
}

func NewAnswerHandler(s *service.AnswerService) *AnswerHandler {
	return &AnswerHandler{service: s}
}

func (h *AnswerHandler) SubmitAnswer(w http.ResponseWriter, r *http.Request) {
	attemptID := r.PathValue("attemptID")

	var input struct {
		QuestionID       string `json:"question_id"`
		SelectedOptionID string `json:"selected_option_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(input.QuestionID) == "" || strings.TrimSpace(input.SelectedOptionID) == "" {
		http.Error(w, "question_id and selected_option_id are required", http.StatusBadRequest)
		return
	}

	answer, err := h.service.SubmitAnswer(r.Context(), attemptID, input.QuestionID, input.SelectedOptionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(answer)
}

func (h *AnswerHandler) UpdateAnswer(w http.ResponseWriter, r *http.Request) {
	attemptID := r.PathValue("attemptID")
	questionID := r.PathValue("questionID")

	var input struct {
		SelectedOptionID string `json:"selected_option_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(input.SelectedOptionID) == "" {
		http.Error(w, "selected_option_id is required", http.StatusBadRequest)
		return
	}

	answer, err := h.service.UpdateAnswer(r.Context(), attemptID, questionID, input.SelectedOptionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(answer)
}
