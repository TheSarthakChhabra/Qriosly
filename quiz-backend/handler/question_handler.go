package handler

import (
	"encoding/json"
	"net/http"
	"quiz-backend/service"
	"strings"
)

type QuestionHandler struct {
	service *service.QuestionService
}

func NewQuestionHandler(s *service.QuestionService) *QuestionHandler {
	return &QuestionHandler{service: s}
}

func (h *QuestionHandler) CreateQuestion(w http.ResponseWriter, r *http.Request) {
	quizID := r.PathValue("id")

	var input struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(input.Text)==""{
		http.Error(w,"text is required", http.StatusBadRequest)
		return
	}

	question, err := h.service.CreateQuestion(r.Context(), quizID, input.Text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(question)
}

func (h *QuestionHandler) GetQuestionsByQuizID(w http.ResponseWriter, r *http.Request) {
	quizID := r.PathValue("id")

	questions, err := h.service.GetQuestionsByQuizID(r.Context(), quizID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questions)
}

func (h *QuestionHandler) GetQuestionByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	question, err := h.service.GetQuestionByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(question)
}
