package repository

import (
	"context"
	"errors"
	"fmt"
	"quiz-backend/model"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AttemptRepository struct {
	db *pgxpool.Pool
}

func NewAttemptRepository(db *pgxpool.Pool) *AttemptRepository {
	return &AttemptRepository{db: db}
}

func (r *AttemptRepository) Save(ctx context.Context, a model.Attempt) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO attempts (id, quiz_id, user_id, started_at, status) VALUES ($1, $2, $3, $4, $5)`,
		a.ID, a.QuizID, a.UserID, a.StartedAt, a.Status,
	)
	if err != nil {
		return fmt.Errorf("failed to save attempt: %w", err)
	}
	return nil
}

func (r *AttemptRepository) FindByID(ctx context.Context, id string) (model.Attempt, error) {
	var a model.Attempt
	err := r.db.QueryRow(ctx,
		`SELECT id, quiz_id, user_id, started_at, status, submitted_at, score FROM attempts WHERE id = $1`, id,
	).Scan(&a.ID, &a.QuizID, &a.UserID, &a.StartedAt, &a.Status, &a.SubmittedAt, &a.Score)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Attempt{}, errors.New("attempt not found")
		}
		return model.Attempt{}, fmt.Errorf("failed to find attempt: %w", err)
	}
	return a, nil
}

func (r *AttemptRepository) MarkSubmitted(ctx context.Context, id, status string, submittedAt time.Time, score int) error {
	_, err := r.db.Exec(ctx,
		`UPDATE attempts SET status = $1, submitted_at = $2, score = $3 WHERE ID = $4`,
		status, submittedAt, score, id,
	)
	if err != nil {
		return fmt.Errorf("failed to mark attempt submitted: %w", err)
	}
	return nil
}

func (r *AnswerRepository) FindByAttemptID(ctx context.Context, attemptID string) ([]model.Answer, error) {
	rows, err := r.db.Query(ctx,
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

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return answers, nil
}
