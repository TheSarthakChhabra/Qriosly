package handler

import (
	"encoding/json"
	"net/http"

	"quiz-backend/apperror"
	"quiz-backend/service"
)

type AdminHandler struct {
	service *service.AdminService
}

func NewAdminHandler(s *service.AdminService) *AdminHandler {
	return &AdminHandler{service: s}
}

func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.ListUsers(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, users)
}

func (h *AdminHandler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	targetUserID := r.PathValue("id")

	var input struct {
		Role string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, apperror.BadRequest("INVALID_BODY", "invalid request body"))
		return
	}

	requesterID, ok := UserIDFromContext(r.Context())
	if !ok {
		writeError(w, apperror.Unauthorized("UNAUTHENTICATED", "authentication required"))
		return
	}

	user, err := h.service.UpdateUserRole(r.Context(), requesterID, targetUserID, input.Role)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, user)
}
