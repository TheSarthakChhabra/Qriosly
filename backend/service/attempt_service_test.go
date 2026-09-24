package service

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"quiz-backend/model"
	"testing"
	"time"
)

func setupScoringFixture() (*fakeAttemptRepo, *fakeQuizRepo, *fakeAnswerRepo, *fakeOptionRepo, model.Attempt) {
	quizRepo := newFakeQuizRepo()
	quizRepo.quizzes["quiz-1"] = model.Quiz{ID: "quiz-1", DurationMinutes: 10}

	optionRepo := newFakeOptionRepo()
	optionRepo.options["opt-correct"] = model.Option{ID: "opt-correct", QuestionID: "q1", IsCorrect: true}
	optionRepo.options["opt-wrong"] = model.Option{ID: "opt-wrong", QuestionID: "q1", IsCorrect: false}
	optionRepo.options["opt-correct-2"] = model.Option{ID: "opt-correct-2", QuestionID: "q2", IsCorrect: true}
	optionRepo.options["opt-wrong-2"] = model.Option{ID: "opt-wrong-2", QuestionID: "q2", IsCorrect: false}

	attemptRepo := newFakeAttemptRepo()
	attempt := model.Attempt{ID: "attempt-1", QuizID: "quiz-1", UserID: "student-1", StartedAt: time.Now(), Status: StatusInProgress}
	attemptRepo.attempts[attempt.ID] = attempt

	return attemptRepo, quizRepo, newFakeAnswerRepo(), optionRepo, attempt
}

func TestStartAttempt_Success(t *testing.T) {
	quizRepo := newFakeQuizRepo()
	quizRepo.quizzes["quiz-1"] = model.Quiz{ID: "quiz-1", DurationMinutes: 10}
	svc := NewAttemptService(newFakeAttemptRepo(), quizRepo, newFakeAnswerRepo(), newFakeOptionRepo(), (*pgxpool.Pool)(nil))

	attempt, err := svc.StartAttempt(context.Background(), "quiz-1", "student-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if attempt.Status != StatusInProgress {
		t.Errorf("expected status %q, got %q", StatusInProgress, attempt.Status)
	}
}

/* 
NOTE: SubmitAttempt is not covered by these unit tests. As of step 16, Submit Attempt runs its critical section inside a real PostgreSQL
	transaction (repository.RunInTx), which requires a genuine *pgxpool. Pool to call Begin() on. The in-memory fakes used elsewhere in 
	this file have no way to simulate a real transaction or row lock, so testing this method here would either panic (nil pool) or only 
	test that a fake obediently does what it's told — not that Postgres actually serializes concurrentsubmissions. SubmitAttempt's scoring 
	logic, status/deadline checks, and concurrency guarantee are instead verified by:
		1. Manual/integration testing against a real database (see: concurrent curl test, two simultaneous POST /attempts/{id}/submit)
 		2. A future dedicated integration test suite that spins up a real or containerized Postgres instance, if this project grows to need one.
*/ 	 

func TestGetAttemptByID_OtherStudent_Forbidden(t *testing.T) {
	attemptRepo, quizRepo, answerRepo, optionRepo, attempt := setupScoringFixture()
	svc := NewAttemptService(attemptRepo, quizRepo, answerRepo, optionRepo, (*pgxpool.Pool)(nil))
	_, err := svc.GetAttemptByID(context.Background(), attempt.ID, "some-other-student", model.RoleStudent)

	if err == nil {
		t.Fatal("expected a forbidden error, got nil")
	}
}

func TestGetAttemptsForQuiz_OtherTeacher_Forbiddent(t *testing.T) {
	quizRepo := newFakeQuizRepo()
	quizRepo.quizzes["quiz-1"] = model.Quiz{ID: "quiz-1", DurationMinutes: 10, CreatedBy: "teacher-owner"}
	svc := NewAttemptService(newFakeAttemptRepo(), quizRepo, newFakeAnswerRepo(), newFakeOptionRepo(), (*pgxpool.Pool)(nil))
	_, err := svc.GetAttemptsForQuiz(context.Background(), "quiz-1", "teacher-someone-else", model.RoleTeacher)

	if err == nil {
		t.Fatal("expected a forbidden error, go nil")
	}
}
