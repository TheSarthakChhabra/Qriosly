package repository
import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5"
	"quiz-backend/model"
	"errors"
)

type OptionRepository struct {
	db *pgxpool.Pool
}

func NewOptionRepository(db *pgxpool.Pool) *OptionRepository {
	return &OptionRepository{db: db}
}

func (r *OptionRepository) Save(ctx context.Context, o model.Option) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO options (id, question_id, text, is_correct) VALUES($1, $2, $3, $4)`,
		o.ID, o.QuestionID, o.Text, o.IsCorrect,
	)
	if err != nil {
		return fmt.Errorf("failed to save option: %w", err)
	}
	return nil
}

func (r *OptionRepository) FindByQuestionID(ctx context.Context, questionID string) ([]model.Option, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, question_id, text, is_correct FROM options WHERE question_id=$1`, questionID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query options: %w", err)
	}
	defer rows.Close()
	var options []model.Option
	for rows.Next() {
		var o model.Option
		if err := rows.Scan(&o.ID, &o.QuestionID, &o.Text, &o.IsCorrect); err != nil {
			return nil, fmt.Errorf("failed to scan option row: %w", err)
		}
		options = append(options, o)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}
	return options, nil
}

func (r *OptionRepository) FindByID(ctx context.Context, id string) (model.Option, error) {
	var o model.Option
	err := r.db.QueryRow(ctx,
		`SELECT id, question_id, text, is_correct FROM options WHERE id = $1`, id,
	).Scan(&o.ID, &o.QuestionID, &o.Text, &o.IsCorrect)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Option{}, errors.New("option not found")
		}
		return model.Option{}, fmt.Errorf("failed to find option: %w", err)
	}

	return o, nil
}