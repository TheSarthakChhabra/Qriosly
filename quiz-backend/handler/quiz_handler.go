package handler

import (
	"encoding/json"
	"net/http"
	"quiz-backend/apperror"
	"quiz-backend/service"
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
		writeError(w, apperror.BadRequest("INVALID_BODY", "invalid request body"))
		return
	}

	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		writeError(w, apperror.Unauthorized("UNAUTHENTICATED", "authentication required"))
		return
	}

	quiz, err := h.service.CreateQuiz(r.Context(), input.Title, input.Description, input.DurationMinutes, userID)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, quiz)
}
