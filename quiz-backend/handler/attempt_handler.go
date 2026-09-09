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

func (h *AttemptHandler) GetAttemptByID(w http.ResponseWriter, r *http.Request) {
	attemptID := r.PathValue("attemptID")
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "uanauthorised", http.StatusUnauthorized)
		return
	}

	role, _ := RoleFromContext(r.Context())
	attempt, err := h.service.GetAttemptByID(r.Context(), attemptID, userID, role)
	if err != nil {
		if err.Error() == "forbidden: you do not have access to this attempt" {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attempt)
}

func (h *AttemptHandler) GetMyAttempts(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	attempts, err := h.service.GetMyAttempts(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attempts)
}

func (h *AttemptHandler) GetAttemptsForQuiz(w http.ResponseWriter, r *http.Request) {
	quizID := r.PathValue("quizID")
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	role, _ := RoleFromContext(r.Context())
	summaries, err := h.service.GetAttemptsForQuiz(r.Context(), quizID, userID, role)
	if err != nil {
		if err.Error() == "forbidden: you do not own this quiz" {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summaries)
}
