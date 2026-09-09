package service

import (
	"context"
	"errors"
	"quiz-backend/model"
	"strings"

	"github.com/google/uuid"
)

type QuizService struct {
	repo QuizRepo
}

func NewQuizService(repo QuizRepo) *QuizService {
	return &QuizService{repo: repo}
}

func (s *QuizService) CreateQuiz(ctx context.Context, title, description string, durationMinutes int, createdBy string) (model.Quiz, error) {
	if strings.TrimSpace(title) ==""{
		return model.Quiz{}, errors.New("title is required")
	}
	if durationMinutes <= 0{
		return model.Quiz{}, errors.New("duration must be positive")
	}
	q := model.Quiz{
		ID:              uuid.NewString(),
		Title:           title,
		Description:     description,
		DurationMinutes: durationMinutes,
		CreatedBy:       createdBy,
	}
	if err := s.repo.Save(ctx, q); err != nil {
		return model.Quiz{}, err
	}
	return q, nil
}

func (s *QuizService) GetAllQuizzes(ctx context.Context) ([]model.Quiz, error) {
	return s.repo.FindAll(ctx)
}

func (s *QuizService) GetQuizByID(ctx context.Context, id string) (model.Quiz, error) {
	return s.repo.FindByID(ctx, id)
}
