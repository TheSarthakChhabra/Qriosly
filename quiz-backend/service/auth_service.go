package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"quiz-backend/apperror"
	"quiz-backend/model"
	"time"
)

type AuthService struct {
	userRepo  UserRepo
	jwtSecret []byte
}

func NewAuthService(userRepo UserRepo, jwtSecret string) *AuthService {
	return &AuthService{userRepo: userRepo, jwtSecret: []byte(jwtSecret)}
}

func (s *AuthService) Register(ctx context.Context, name, email, password, role string) (model.User, error) {
	if _, err := s.userRepo.FindByEmail(ctx, email); err == nil {
		return model.User{}, apperror.Conflict("EMAIL_ALREADY_EXISTS", "A user with this email already exists")
	}

	if role == "" {
		role = model.RoleStudent
	}

	if role != model.RoleAdmin && role != model.RoleStudent && role != model.RoleTeacher {
		return model.User{}, apperror.BadRequest("INVALID_ROLE", "Invalid role")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	u := model.User{
		ID:           uuid.NewString(),
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		Role:         role,
		CreatedAt:    time.Now(),
	}
	if err := s.userRepo.Save(ctx, u); err != nil {
		return model.User{}, err
	}
	return u, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	u, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", apperror.Unauthorized("INVALID_CREDENTIALS", "invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", apperror.Unauthorized("INVALID_CREDENTIALS", "invalid email or password")
	}

	token, err := s.generateToken(u.ID, u.Role)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return token, nil
}

func (s *AuthService) generateToken(userID, role string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func ParseToken(tokenString, secret string) (userID, role string, err error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return "", "", errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("invalid token claims")
	}

	userID, ok = claims["sub"].(string)
	if !ok {
		return "", "", errors.New("invalid token subject")
	}

	role, ok = claims["role"].(string)
	if !ok {
		return "", "", errors.New("invalid token role")
	}

	return userID, role, nil
}
