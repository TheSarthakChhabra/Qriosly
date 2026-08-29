package repository

import (
	"context"
	"errors"
	"fmt"
	"quiz-backend/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type QuizRepository struct {
	db *pgxpool.Pool
}

func NewQuizRepository(db *pgxpool.Pool) *QuizRepository {
	return &QuizRepository{db: db}
}

func (r *QuizRepository) Save(ctx context.Context, q model.Quiz) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO quizzes (id, title, description) VALUES($1, $2, $3)`,
		q.ID, q.Title, q.Description,
	)
	if err != nil {
		return fmt.Errorf("failed to save quiz: %w", err)
	}
	return nil
}

func (r *QuizRepository) FindAll(ctx context.Context) ([]model.Quiz, error) {
	rows, err := r.db.Query(ctx, `SELECT id, title, description FROM quizzes`)
	if err != nil {
		return nil, fmt.Errorf("failed to query quizzes: %w", err)
	}
	defer rows.Close()
	var quizzes []model.Quiz
	for rows.Next() {
		var q model.Quiz
		if err := rows.Scan(&q.ID, &q.Title, &q.Description); err != nil {
			return nil, fmt.Errorf("failed to scan quiz: %w", err)
		}
		quizzes = append(quizzes, q)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}
	return quizzes, nil
}

func (r *QuizRepository) FindByID(ctx context.Context, id string) (model.Quiz, error) {
	var q model.Quiz
	err := r.db.QueryRow(ctx,
		`SELECT id, title, description FROM quizzes WHERE id = $1`, id,
	).Scan(&q.ID, &q.Title, &q.Description)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Quiz{}, errors.New("quiz not found")
		}
		return model.Quiz{}, fmt.Errorf("failed to find quiz: %w", err)
	}
	return q, nil
}
