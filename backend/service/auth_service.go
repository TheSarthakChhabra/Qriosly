package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"quiz-backend/apperror"
	"quiz-backend/model"
	"time"
	"unicode"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  UserRepo
	jwtSecret []byte
}

func NewAuthService(userRepo UserRepo, jwtSecret string) *AuthService {
	return &AuthService{userRepo: userRepo, jwtSecret: []byte(jwtSecret)}
}

func isValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasLetter, hasDigit := false, false
	for _, c := range password {
		switch {
		case unicode.IsLetter(c):
			hasLetter = true
		case unicode.IsDigit(c):
			hasDigit = true
		}
	}
	return hasLetter && hasDigit
}

func (s *AuthService) Register(ctx context.Context, name, email, password, role string) (model.User, error) {
	if !isValidPassword(password) {
		return model.User{}, apperror.BadRequest("WEAK_PASSWORD", "password must be at least 8 characters and include both letters and numbers")
	}
	if _, err := s.userRepo.FindByEmail(ctx, email); err == nil {
		return model.User{}, apperror.Conflict("EMAIL_ALREADY_EXISTS", "A user with this email already exists")
	}
	role = model.RoleStudent
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
	slog.Info("user logged in", "user_id", u.ID, "role", u.Role)
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
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{"HS256"}))
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
