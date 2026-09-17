package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"quiz-backend/apperror"
)

type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type errorResponse struct {
	Error errorBody `json:"error"`
}

func writeError(w http.ResponseWriter, err error) {
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		writeJSON(w, appErr.Status, errorResponse{
			Error: errorBody{Code: appErr.Code, Message: appErr.Message},
		})
		return
	}
	log.Printf("internal error: %v", err)
	writeJSON(w, http.StatusInternalServerError, errorResponse{
		Error: errorBody{Code: "INTERNAL_ERROR", Message: "An unexpected error occured"},
	})
}

func writeJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}
