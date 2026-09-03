package service

import (
	"context"
	"quiz-backend/model"
	"quiz-backend/repository"

	"github.com/google/uuid"
)

type QuizService struct {
	repo *repository.QuizRepository
}

func NewQuizService(repo *repository.QuizRepository) *QuizService {
	return &QuizService{repo: repo}
}

func (s *QuizService) CreateQuiz(ctx context.Context, title, description string, durationMinutes int) (model.Quiz, error) {
	q := model.Quiz{
		ID:          uuid.NewString(),
		Title:       title,
		Description: description,
		DurationMinutes: durationMinutes,
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
