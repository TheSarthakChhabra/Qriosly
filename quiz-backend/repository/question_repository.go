package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"quiz-backend/model"
)

type QuestionRepository struct {
	db *pgxpool.Pool
}

func NewQuestionRepository(db *pgxpool.Pool) *QuestionRepository {
	return &QuestionRepository{db: db}
}

func (r *QuestionRepository) Save(ctx context.Context, q model.Question) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO questions (id, quiz_id, text) VALUES ($1, $2, $3)`,
		q.ID, q.QuizID, q.Text,
	)
	if err != nil {
		return fmt.Errorf("failed to save the question: %w", err)
	}
	return nil
}

func (r *QuestionRepository) FindByQuizID(ctx context.Context, quizID string) ([]model.Question, error) {
	// finding all questions of the quiz_id
	rows, err := r.db.Query(ctx,
		`SELECT id, quiz_id, text FROM questions WHERE quiz_id = $1`, quizID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query questions: %w", err)
	}
	defer rows.Close()
	var questions []model.Question
	for rows.Next() {
		var q model.Question
		if err := rows.Scan(&q.ID, &q.QuizID, &q.Text); err != nil {
			return nil, fmt.Errorf("failed to scan questions: %w", err)
		}
		questions = append(questions, q)
	}
	return questions, nil
}

func (r *QuestionRepository) FindByID(ctx context.Context, id string) (model.Question, error) {
	// finding a particular question by its id
	var q model.Question
	err := r.db.QueryRow(ctx,
		`SELECT id, quiz_id, text FROM questions WHERE id = $1`, id,
	).Scan(&q.ID, &q.QuizID, &q.Text)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Question{}, errors.New("question not found")
		}
		return model.Question{}, fmt.Errorf("failed to find the question:%w", err)
	}
	return q, nil
}
