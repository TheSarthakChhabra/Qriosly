package service

import (
	"context"
	"errors"
	"log/slog"
	"quiz-backend/apperror"
	"quiz-backend/model"
	"quiz-backend/repository"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
	pool        *pgxpool.Pool
}

func NewAttemptService(attemptRepo AttemptRepo, quizRepo QuizRepo, answerRepo AnswerRepo, optionRepo OptionRepo, pool *pgxpool.Pool) *AttemptService {
	return &AttemptService{
		attemptRepo: attemptRepo, quizRepo: quizRepo, answerRepo: answerRepo, optionRepo: optionRepo, pool: pool,
	}
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

	var result model.Attempt
	txErr := repository.RunInTx(ctx, s.pool, func(tx pgx.Tx) error {
		attempt, err := s.attemptRepo.FindByIDForUpdate(ctx, tx, attemptID)
		if err != nil {
			return apperror.NotFound("ATTEMPT_NOT_FOUND", "attempt not found")
		}
		if attempt.Status != StatusInProgress {
			return apperror.Conflict("ATTEMPT_ALREADY_SUBMITTED", "This attempt has already been submitted")
		}

		quiz, err := s.quizRepo.FindByID(ctx, attempt.QuizID)
		if err != nil {
			return apperror.NotFound("QUIZ_NOT_FOUND", "Quiz not found")
		}

		deadline := attempt.StartedAt.Add(time.Duration(quiz.DurationMinutes) * time.Minute)
		if time.Now().After(deadline) {
			return apperror.BadRequest("ATTEMPT_EXPIRED", "The time limit for this attempt has passed")
		}

		answers, err := s.answerRepo.FindByAttemptIDTx(ctx, tx, attemptID)
		if err != nil {
			return apperror.Internal("ANSWERS_FETCH_FAILED", "failed to fetch answers")
		}
		score := 0
		for _, ans := range answers {
			option, err := s.optionRepo.FindByIDTx(ctx, tx, ans.SelectedOptionID)
			if err != nil {
				continue
			}
			if option.IsCorrect {
				score++
			}
		}
		submittedAt := time.Now()
		if err := s.attemptRepo.MarkSubmittedTx(ctx, tx, attemptID, StatusSubmitted, submittedAt, score); err != nil {
			return apperror.Internal("ATTEMPT_UPDATE_FAILED", "failed to mark attempt submitted")
		}
		attempt.Status = StatusSubmitted
		attempt.SubmittedAt = &submittedAt
		attempt.Score = &score
		result = attempt
		return nil
	})
	if txErr != nil {
		return model.Attempt{}, txErr
	}
	slog.Info("attempt submitted", "attempt_id", result.ID, "user_id", result.UserID, "score", *result.Score)
	return result, nil
}

func (s *AttemptService) GetAttemptByID(ctx context.Context, attemptID, requesterID, requesterRole string) (model.Attempt, error) {
	attempt, err := s.attemptRepo.FindByID(ctx, attemptID)
	if err != nil {
		return model.Attempt{}, apperror.NotFound("ATTEMPT_NOT_FOUND", "Attempt not found")
	}

	if requesterRole == model.RoleAdmin || requesterID == attempt.UserID{
		return attempt, nil
	}

	if requesterRole == model.RoleTeacher {
		quiz, err := s.quizRepo.FindByID(ctx, attempt.QuizID)
		if err == nil && quiz.CreatedBy == requesterID {
			return attempt, nil
		}
	}
	return model.Attempt{}, apperror.Forbidden("ATTEMPT_ACCESS_DENIED", "you donot have access to this attempt")
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
