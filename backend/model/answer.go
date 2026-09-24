package model

type Answer struct {
	ID               string `json:"id"`
	AttemptID        string `json:"attempt_id"`
	QuestionID       string `json:"question_id"`
	SelectedOptionID string `json:"selected_option_id"`
}