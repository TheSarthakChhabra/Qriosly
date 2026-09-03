package handler

import (
	"encoding/json"
	"net/http"

	"quiz-backend/service"
)

type AttemptHandler struct {
	service *service.AttemptService
}

func NewAttemptHandler(s *service.AttemptService) *AttemptHandler {
	return &AttemptHandler{service: s}
}

func (h *AttemptHandler) StartAttempt(w http.ResponseWriter, r *http.Request) {
	quizID := r.PathValue("id")

	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	attempt, err := h.service.StartAttempt(r.Context(), quizID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(attempt)
}

func (h *AttemptHandler) SubmitAttempt(w http.ResponseWriter, r *http.Request) {
	attemptID := r.PathValue("attemptID")
	attempt, err := h.service.SubmitAttempt(r.Context(), attemptID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attempt)
}
