package model

import "time"

type AttemptSummary struct {
	AttemptID   string     `json:"attempt_id"`
	QuizID      string     `json:"quiz_id"`
	UserID      string     `json:"user_id"`
	UserName    string     `json:"user_name"`
	UserEmail   string     `json:"user_email"`
	StartedAt   time.Time  `json:"started_at"`
	Status      string     `json:"status"`
	SubmittedAt *time.Time `json:"submitted_at,omitempty"`
	Score       *int       `json:"score,omitempty"`
}
