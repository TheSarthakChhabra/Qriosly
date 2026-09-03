package handler

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"quiz-backend/service"
	"strings"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(input.Name) ==""{
		http.Error(w,"name is required",http.StatusBadRequest)
		return
	}

	if _,err := mail.ParseAddress(input.Email); err!=nil{
		http.Error(w,"invalid email address", http.StatusBadRequest)
		return
	}

	if len(input.Password)<8{
		http.Error(w, "password must be atleast 8 characters", http.StatusBadRequest)
		return
	}

	user, err := h.service.Register(r.Context(), input.Name, input.Email, input.Password, input.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Cpntent-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	token, err := h.service.Login(r.Context(), input.Email, input.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
