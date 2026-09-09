package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"quiz-backend/model"
	"time"
)

const mockUserID = "00000000-0000-0000-0000-000000000001"

const (
	StatusInProgress = "in_progress"
	StatusSubmitted  = "submitted"
)

type AttemptService struct {
	attemptRepo AttemptRepo
	quizRepo    QuizRepo
	answerRepo  AnswerRepo
	optionRepo  OptionRepo
}

func NewAttemptService(attemptRepo AttemptRepo, quizRepo QuizRepo, answerRepo AnswerRepo, optionRepo OptionRepo) *AttemptService {
	return &AttemptService{attemptRepo: attemptRepo, quizRepo: quizRepo, answerRepo: answerRepo, optionRepo: optionRepo}
}

func (s *AttemptService) StartAttempt(ctx context.Context, quizID, userID string) (model.Attempt, error) {
	if _, err := s.quizRepo.FindByID(ctx, quizID); err != nil {
		return model.Attempt{}, errors.New("cannot start attempt: quiz not found")
	}

	a := model.Attempt{
		ID:        uuid.NewString(),
		QuizID:    quizID,
		UserID:    userID,
		StartedAt: time.Now(),
		Status:    StatusInProgress,
	}
	if err := s.attemptRepo.Save(ctx, a); err != nil {
		return model.Attempt{}, err
	}
	return a, nil
}

func (s *AttemptService) SubmitAttempt(ctx context.Context, attemptID string) (model.Attempt, error) {

	attempt, err := s.attemptRepo.FindByID(ctx, attemptID)
	if err != nil {
		return model.Attempt{}, errors.New("cannot submit the attempt: attempt not found")
	}
	if attempt.Status != StatusInProgress {
		return model.Attempt{}, errors.New("cannot submit the attempt: attempt not in progress")
	}

	quiz, err := s.quizRepo.FindByID(ctx, attempt.QuizID)
	if err != nil {
		return model.Attempt{}, errors.New("cannot submit attempt: quiz not found")
	}

	deadline := attempt.StartedAt.Add(time.Duration(quiz.DurationMinutes) * time.Minute)
	if time.Now().After(deadline) {
		return model.Attempt{}, errors.New("cannot submit attempt: time limit exceeded")
	}

	answers, err := s.answerRepo.FindByAttemptID(ctx, attemptID)
	if err != nil {
		return model.Attempt{}, err
	}
	score := 0
	for _, ans := range answers {
		option, err := s.optionRepo.FindByID(ctx, ans.SelectedOptionID)
		if err != nil {
			continue
		}
		if option.IsCorrect {
			score++
		}
	}
	submittedAt := time.Now()
	if err := s.attemptRepo.MarkSubmitted(ctx, attemptID, StatusSubmitted, submittedAt, score); err != nil {
		return model.Attempt{}, err
	}
	attempt.Status = StatusSubmitted
	attempt.SubmittedAt = &submittedAt
	attempt.Score = &score
	return attempt, nil
}

func (s *AttemptService) GetAttemptByID(ctx context.Context, attemptID, requesterID, requesterRole string) (model.Attempt, error) {
	attempt, err := s.attemptRepo.FindByID(ctx, attemptID)
	if err != nil {
		return model.Attempt{}, errors.New("attempt not found")
	}

	if requesterRole == model.RoleAdmin {
		return attempt, nil
	}
	if requesterID == attempt.UserID {
		return attempt, nil
	}
	if requesterRole == model.RoleTeacher {
		quiz, err := s.quizRepo.FindByID(ctx, attempt.QuizID)
		if err == nil && quiz.CreatedBy == requesterID {
			return attempt, nil
		}
	}
	return model.Attempt{}, errors.New("forbidden: you do not have access to this attempt")
}

func (s *AttemptService) GetMyAttempts(ctx context.Context, userID string) ([]model.Attempt, error) {
	return s.attemptRepo.FindByUserID(ctx, userID)
}

func (s *AttemptService) GetAttemptsForQuiz(ctx context.Context, quizID, requesterID, requesterRole string) ([]model.AttemptSummary, error) {
	quiz, err := s.quizRepo.FindByID(ctx, quizID)
	if err != nil {
		return nil, errors.New("quiz not found")
	}

	if requesterRole != model.RoleAdmin && quiz.CreatedBy != requesterID {
		return nil, errors.New("forbidden: you don not own this quiz")
	}
	return s.attemptRepo.FindSummariesByQuizID(ctx, quizID)
}
