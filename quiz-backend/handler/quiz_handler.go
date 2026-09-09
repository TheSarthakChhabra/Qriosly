package handler

import (
	"encoding/json"
	"net/http"
	"quiz-backend/service"
	"strings"
)

type QuizHandler struct {
	service *service.QuizService
}

func NewQuizHandler(s *service.QuizService) *QuizHandler {
	return &QuizHandler{service: s}
}

func (h *QuizHandler) GetAllQuizzes(w http.ResponseWriter, r *http.Request) {
	quizzes, err := h.service.GetAllQuizzes(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quizzes)
}

func (h *QuizHandler) GetQuizByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	quiz, err := h.service.GetQuizByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quiz)
}

func (h *QuizHandler) CreateQuiz(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title           string `json:"title"`
		Description     string `json:"description"`
		DurationMinutes int    `json:"duration_minutes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(input.Title) == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	if input.DurationMinutes <= 0 {
		http.Error(w, "duration_minutes must be positive", http.StatusBadRequest)
		return
	}

	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	quiz, err := h.service.CreateQuiz(r.Context(), input.Title, input.Description, input.DurationMinutes, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(quiz)
}
