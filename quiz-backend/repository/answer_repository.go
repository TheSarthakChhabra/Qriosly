package repository

import (
	"context"
	"fmt"
	"errors"
	"quiz-backend/model"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5"
)

type AnswerRepository struct {
	db *pgxpool.Pool
}

func NewAnswerRepository(db *pgxpool.Pool) *AnswerRepository {
	return &AnswerRepository{db: db}
}

func (r *AnswerRepository) Save(ctx context.Context, a model.Answer) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO answers (id, attempt_id, question_id, selected_option_id) VALUES ($1, $2, $3, $4)`,
		a.ID, a.AttemptID, a.QuestionID, a.SelectedOptionID,
	)
	if err != nil {
		return fmt.Errorf("failed to save answer: %w", err)
	}
	return nil
}

func (r *AnswerRepository) FindByAttemptAndQuestion(ctx context.Context, attemptID, questionID string) (model.Answer, error) {
	var a model.Answer
	err := r.db.QueryRow(ctx,
		`SELECT id, attempt_id, question_id, selected_option_id FROM answers WHERE attempt_id = $1 AND question_id = $2`,
		attemptID, questionID,
	).Scan(&a.ID, &a.AttemptID, &a.QuestionID, &a.SelectedOptionID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Answer{}, errors.New("answer not found")
		}
		return model.Answer{}, fmt.Errorf("failed to find answer: %w", err)
	}

	return a, nil
}

func (r *AnswerRepository) UpdateSelectedOption(ctx context.Context, answerID, selectedOptionID string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE answers SET selected_option_id = $1 WHERE id = $2`,
		selectedOptionID, answerID,
	)
	if err != nil {
		return fmt.Errorf("failed to update answer: %w", err)
	}
	return nil
}