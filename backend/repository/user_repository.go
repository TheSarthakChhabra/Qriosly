package repository

import (
	"context"
	"errors"
	"fmt"
	"quiz-backend/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Save(ctx context.Context, u model.User) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO users(id,name,email,password_hash,role,created_at) VALUES($1,$2,$3,$4,$5,$6)`,
		u.ID, u.Name, u.Email, u.PasswordHash, u.Role, u.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}
	return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (model.User, error) {
	var u model.User
	err := r.db.QueryRow(ctx,
		`SELECT id, name, email, password_hash, role, created_at FROM users WHERE email = $1`, email,
	).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, errors.New("user not found")
		}
		return model.User{}, fmt.Errorf("failed to find user by email: %w", err)
	}
	return u, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (model.User, error) {
	var u model.User
	err := r.db.QueryRow(ctx,
		`SELECT id, name, email, password_hash, role, created_at FROM users WHERE id=$1`, id,
	).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, errors.New("user not found")
		}
		return model.User{}, fmt.Errorf("failed to find user by id: %w", err)
	}
	return u, nil
}

func (r *UserRepository) FindAll(ctx context.Context) ([]model.User, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name, email, password_hash, role, created_at FROM users`)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	users := []model.User{}
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}
	return users, nil
}

func (r *UserRepository) UpdateRole(ctx context.Context, id, role string) error {
	_, err := r.db.Exec(ctx, `UPDATE users SET role = $1 WHERE id = $2`, role, id)
	if err != nil {
		return fmt.Errorf("failed to update user role: %w", err)
	}
	return nil
}