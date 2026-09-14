package service

import (
	"context"
	"errors"
	"quiz-backend/model"
	"github.com/google/uuid"
)

type AnswerService struct {
	answerRepo   AnswerRepo
	attemptRepo  AttemptRepo
	questionRepo QuestionRepo
	optionRepo   OptionRepo
}

func NewAnswerService(
	answerRepo AnswerRepo,
	attemptRepo AttemptRepo,
	questionRepo QuestionRepo,
	optionRepo OptionRepo,
) *AnswerService {
	return &AnswerService{
		answerRepo:   answerRepo,
		attemptRepo:  attemptRepo,
		questionRepo: questionRepo,
		optionRepo:   optionRepo,
	}
}

func (s *AnswerService) SubmitAnswer(ctx context.Context, attemptID, questionID, selectedOptionID string) (model.Answer, error) {
	attempt, err := s.attemptRepo.FindByID(ctx, attemptID)
	if err != nil {
		return model.Answer{}, errors.New("cannot submit answer: attempt not found")
	}

	question, err := s.questionRepo.FindByID(ctx, questionID)
	if err != nil {
		return model.Answer{}, errors.New("cannot submit answer: question not found")
	}
	if question.QuizID != attempt.QuizID {
		return model.Answer{}, errors.New("cannot submit answer: question does not belong to this attempt's quiz")
	}

	option, err := s.optionRepo.FindByID(ctx, selectedOptionID)
	if err != nil {
		return model.Answer{}, errors.New("cannot submit answer: selected option not found")
	}
	if option.QuestionID != questionID {
		return model.Answer{}, errors.New("cannot submit answer: selected option does not belong to this question")
	}

	a := model.Answer{
		ID:               uuid.NewString(),
		AttemptID:        attemptID,
		QuestionID:       questionID,
		SelectedOptionID: selectedOptionID,
	}
	if err := s.answerRepo.Save(ctx, a); err != nil {
		return model.Answer{}, err
	}
	return a, nil
}

func (s *AnswerService) UpdateAnswer(ctx context.Context, attemptID, questionID, selectedOptionID string) (model.Answer, error) {
	attempt, err := s.attemptRepo.FindByID(ctx, attemptID)
	if err != nil {
		return model.Answer{}, errors.New("cannot update answer: attempt not found")
	}
	if attempt.Status != StatusInProgress {
		return model.Answer{}, errors.New("cannot update answer: attempt is not in progress")
	}

	question, err := s.questionRepo.FindByID(ctx, questionID)
	if err != nil {
		return model.Answer{}, errors.New("cannot update answer: question not found")
	}
	if question.QuizID != attempt.QuizID {
		return model.Answer{}, errors.New("cannot update answer: question does not belong to this attempt's quiz")
	}

	option, err := s.optionRepo.FindByID(ctx, selectedOptionID)
	if err != nil {
		return model.Answer{}, errors.New("cannot update answer: selected option not found")
	}
	if option.QuestionID != questionID {
		return model.Answer{}, errors.New("cannot update answer: selected option does not belong to this question")
	}

	existing, err := s.answerRepo.FindByAttemptAndQuestion(ctx, attemptID, questionID)
	if err != nil {
		return model.Answer{}, errors.New("cannot update answer: no existing answer for this question")
	}

	if err := s.answerRepo.UpdateSelectedOption(ctx, existing.ID, selectedOptionID); err != nil {
		return model.Answer{}, err
	}

	existing.SelectedOptionID = selectedOptionID
	return existing, nil
}