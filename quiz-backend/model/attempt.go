package model

import "time"

type Attempt struct {
	ID          string     `json:"id"`
	QuizID      string     `json:"quiz_id"`
	UserID      string     `json:"user_id"`
	StartedAt   time.Time  `json:"started_at"`
	Status      string     `json:"status"`
	SubmittedAt *time.Time `json:"submitted_at,omitempty"`
	Score       *int       `json:"score,omitempty"`
}
