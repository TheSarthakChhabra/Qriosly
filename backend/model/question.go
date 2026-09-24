package model

type Question struct {
	ID     string `json:"id"`
	QuizID string `json:"quiz_id"`
	Text   string `json:"text"`
}
