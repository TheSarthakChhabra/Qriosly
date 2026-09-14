package repository

import (
	"context"
	"errors"
	"fmt"
	"quiz-backend/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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

func (r *AnswerRepository) FindByAttemptIDTx(ctx context.Context, tx pgx.Tx, attemptID string) ([]model.Answer, error) {
	rows, err := tx.Query(ctx,
		`SELECT id, attempt_id, question_id, selected_option_id FROM answers WHERE attempt_id = $1`, attemptID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query answers: %w", err)
	}
	defer rows.Close()

	var answers []model.Answer
	for rows.Next() {
		var a model.Answer
		if err := rows.Scan(&a.ID, &a.AttemptID, &a.QuestionID, &a.SelectedOptionID); err != nil {
			return nil, fmt.Errorf("failed to scan answer row: %w", err)
		}
		answers = append(answers, a)
	}
	return answers, rows.Err()
}
