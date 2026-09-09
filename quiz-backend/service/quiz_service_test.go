package service

import (
	"context"
	"testing"
)

func TestCreateQuiz_Valid(t *testing.T) {
	repo := newFakeQuizRepo()
	svc := NewQuizService(repo)

	quiz, err := svc.CreateQuiz(context.Background(), "Go Basics", "intro", 10, "teacher-1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if quiz.ID == "" {
		t.Fatalf("expected a generated ID, got empty string")
	}
}

func TestCreateQuiz_Empty_Title_Rejected(t *testing.T) {
	svc := NewQuizService(newFakeQuizRepo())
	_, err := svc.CreateQuiz(context.Background(), "", "intro", 10, "teacher-1")

	if err == nil {
		t.Fatalf("expected an error for empty title, got nil")
	}
}

func TestCreateQuiz_NonPositveDuration_Rejected(t *testing.T) {
	svc := NewQuizService(newFakeQuizRepo())
	_, err := svc.CreateQuiz(context.Background(), "Go Basics", "intro", 0, "teacher-1")
	if err == nil {
		t.Fatalf("expected an error for non-positive duration, got nil")
	}
}
