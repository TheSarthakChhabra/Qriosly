package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"quiz-backend/model"
)

type OptionService struct {
	optionRepo   OptionRepo
	questionRepo QuestionRepo
}

func NewOptionService(optionRepo OptionRepo, questionRepo QuestionRepo) *OptionService {
	return &OptionService{optionRepo: optionRepo, questionRepo: questionRepo}
}

func (s *OptionService) CreateOption(ctx context.Context, questionID, text string, isCorrect bool) (model.Option, error) {
	if _, err := s.questionRepo.FindByID(ctx, questionID); err != nil {
		return model.Option{}, errors.New("cannot create option: question not found")
	}
	o := model.Option{
		ID:         uuid.NewString(),
		QuestionID: questionID,
		Text:       text,
		IsCorrect:  isCorrect,
	}
	if err := s.optionRepo.Save(ctx, o); err != nil {
		return model.Option{}, err
	}
	return o, nil
}

func (s *OptionService) GetOptionsByQuestionID(ctx context.Context, questionID string) ([]model.Option, error) {
	return s.optionRepo.FindByQuestionID(ctx,questionID)
}
