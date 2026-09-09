package service

import (
	"context"
	"quiz-backend/model"
	"time"
)

type QuizRepo interface {
	Save(ctx context.Context, q model.Quiz) error
	FindAll(ctx context.Context) ([]model.Quiz, error)
	FindByID(ctx context.Context, id string) (model.Quiz, error)
}

type QuestionRepo interface {
	Save(ctx context.Context, q model.Question) error
	FindByQuizID(ctx context.Context, quizID string) ([]model.Question, error)
	FindByID(ctx context.Context, questionID string) (model.Question, error)
}

type OptionRepo interface {
	Save(ctx context.Context, o model.Option) error
	FindByQuestionID(ctx context.Context, questionID string) ([]model.Option, error)
	FindByID(ctx context.Context, id string) (model.Option, error)
}

type AttemptRepo interface {
	Save(ctx context.Context, a model.Attempt) error
	FindByID(ctx context.Context, id string) (model.Attempt, error)
	MarkSubmitted(ctx context.Context, id, status string, submittedAt time.Time, score int) error
	FindByUserID(ctx context.Context, userID string) ([]model.Attempt, error)
	FindByQuizID(ctx context.Context, quizID string) ([]model.Attempt, error)
	FindSummariesByQuizID(ctx context.Context, quizID string) ([]model.AttemptSummary, error)
}

type AnswerRepo interface {
	Save(ctx context.Context, a model.Answer) error
	FindByAttemptAndQuestion(ctx context.Context, attemptID, questionID string) (model.Answer, error)
	UpdateSelectedOption(ctx context.Context, answerID, selectedOptionID string) error
	FindByAttemptID(ctx context.Context, attemptID string) ([]model.Answer, error)
}

type UserRepo interface {
	Save(ctx context.Context, u model.User) error
	FindByEmail(ctx context.Context, email string) (model.User, error)
	FindByID(ctx context.Context, id string) (model.User, error)
}
