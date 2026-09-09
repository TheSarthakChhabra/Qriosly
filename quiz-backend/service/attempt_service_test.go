package service

import (
	"context"
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
	svc := NewAttemptService(newFakeAttemptRepo(), quizRepo, newFakeAnswerRepo(), newFakeOptionRepo())

	attempt, err := svc.StartAttempt(context.Background(), "quiz-1", "student-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if attempt.Status != StatusInProgress {
		t.Errorf("expected status %q, got %q", StatusInProgress, attempt.Status)
	}
}

func TestSubmitAttempt_CorrectScore(t *testing.T) {
	attemptRepo, quizRepo, answerRepo, optionRepo, attempt := setupScoringFixture()
	answerRepo.answers["a1"] = model.Answer{ID: "a1", AttemptID: attempt.ID, QuestionID: "q1", SelectedOptionID: "opt-correct"}
	answerRepo.answers["a2"] = model.Answer{ID: "a2", AttemptID: attempt.ID, QuestionID: "q2", SelectedOptionID: "opt-correct-2"}

	svc := NewAttemptService(attemptRepo, quizRepo, answerRepo, optionRepo)
	result, err := svc.SubmitAttempt(context.Background(), attempt.ID)

	if err != nil {
		t.Fatalf("expected no error, found %v", err)
	}
	if result.Score == nil || *result.Score != 2 {
		t.Errorf("expected score 2, got %v", result.Score)
	}
}

func TestSubmitAtteot_IncorrectScore(t *testing.T) {
	attemptRepo, quizRepo, answerRepo, optionRepo, attempt := setupScoringFixture()
	answerRepo.answers["a1"] = model.Answer{ID: "a1", AttemptID: attempt.ID, QuestionID: "q1", SelectedOptionID: "opt-wrong"}
	answerRepo.answers["a2"] = model.Answer{ID: "a2", AttemptID: attempt.ID, QuestionID: "q2", SelectedOptionID: "opt-correct-2"}
	svc := NewAttemptService(attemptRepo, quizRepo, answerRepo, optionRepo)

	result, err := svc.SubmitAttempt(context.Background(), attempt.ID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Score == nil || *result.Score != 1 {
		t.Errorf("expected score 1, got %v", result.Score)
	}
}

func TestSubmitAttempt_AlreadySubmitted_Rejected(t *testing.T) {
	attemptRepo, quizRepo, answerRepo, optionRepo, attempt := setupScoringFixture()
	attempt.Status = StatusSubmitted
	attemptRepo.attempts[attempt.ID] = attempt
	svc := NewAttemptService(attemptRepo, quizRepo, answerRepo, optionRepo)

	_, err := svc.SubmitAttempt(context.Background(), attempt.ID)

	if err == nil {
		t.Fatal("expected an error submitting an already-submitted attempt, got nil")
	}
}

func TestSubmitAttempt_Expired_Rejected(t *testing.T) {
	attemptRepo, quizRepo, answerRepo, optionRepo, attempt := setupScoringFixture()
	attempt.StartedAt = time.Now().Add(-20 * time.Minute)
	attemptRepo.attempts[attempt.ID] = attempt
	svc := NewAttemptService(attemptRepo, quizRepo, answerRepo, optionRepo)

	_, err := svc.SubmitAttempt(context.Background(), attempt.ID)

	if err == nil {
		t.Fatal("expected an error for expired attempt, got nil")
	}
}

func TestGetAttemptByID_OtherStudent_Forbidden(t *testing.T) {
	attemptRepo, quizRepo, answerRepo, optionRepo, attempt := setupScoringFixture()
	svc:=NewAttemptService(attemptRepo, quizRepo, answerRepo, optionRepo)
	_, err:=svc.GetAttemptByID(context.Background(), attempt.ID, "some-other-student", model.RoleStudent)

	if err==nil{
		t.Fatal("expected a forbidden error, got nil")
	}
}

func TestGetAttemptsForQuiz_OtherTeacher_Forbiddent(t *testing.T){
	quizRepo := newFakeQuizRepo()
	quizRepo.quizzes["quiz-1"] = model.Quiz{ID: "quiz-1", DurationMinutes: 10, CreatedBy: "teacher-owner"}
	svc := NewAttemptService(newFakeAttemptRepo(),quizRepo, newFakeAnswerRepo(), newFakeOptionRepo())
	_, err:= svc.GetAttemptsForQuiz(context.Background(), "quiz-1", "teacher-someone-else", model.RoleTeacher)

	if err==nil{
		t.Fatal("expected a forbidden error, go nil")
	}
}