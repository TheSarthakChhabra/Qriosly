package service

import (
	"context"
	"errors"
	"quiz-backend/model"
	"github.com/google/uuid"
)

type QuestionService struct {
	questionRepo QuestionRepo
	quizRepo     QuizRepo
}

func NewQuestionService(questionRepo QuestionRepo, quizRepo QuizRepo) *QuestionService {
	return &QuestionService{questionRepo: questionRepo, quizRepo: quizRepo}
}

func (s *QuestionService) CreateQuestion(ctx context.Context, quizID, text string) (model.Question, error) {
	if _, err := s.quizRepo.FindByID(ctx, quizID); err != nil {
		return model.Question{}, errors.New("cannot create question: quiz not found")
	}
	q := model.Question{
		ID:     uuid.NewString(),
		QuizID: quizID,
		Text:   text,
	}
	if err := s.questionRepo.Save(ctx, q); err != nil {
		return model.Question{}, err
	}
	return q, nil
}

func (s *QuestionService) GetQuestionsByQuizID(ctx context.Context, id string) ([]model.Question, error) {
	return s.questionRepo.FindByQuizID(ctx, id)
}

func (s *QuestionService) GetQuestionByID(ctx context.Context, id string) (model.Question, error) {
	return s.questionRepo.FindByID(ctx, id)
}
