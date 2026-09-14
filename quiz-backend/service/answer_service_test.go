package service

import (
	"context"
	"quiz-backend/model"
	"testing"
	"time"
)

func setupAnswerFixture() (*fakeAnswerRepo, *fakeAttemptRepo, *fakeQuestionRepo, *fakeOptionRepo, model.Attempt) {
	questionRepo := newFakeQuestionRepo()
	questionRepo.questions["q1"] = model.Question{ID: "q1", QuizID: "quiz-1"}
	questionRepo.questions["q-other-quiz"] = model.Question{ID: "q-other-quiz", QuizID: "quiz-2"}

	optionRepo := newFakeOptionRepo()
	optionRepo.options["opt-1"] = model.Option{ID: "opt-1", QuestionID: "q1", IsCorrect: true}
	optionRepo.options["opt-2"] = model.Option{ID: "opt-2", QuestionID: "q1", IsCorrect: false}
	optionRepo.options["opt-other-question"] = model.Option{ID: "opt-other-question", QuestionID: "some-other-question"}

	attemptRepo := newFakeAttemptRepo()
	attempt := model.Attempt{ID: "attempt-1", QuizID: "quiz-1", UserID: "student-1", StartedAt: time.Now(), Status: StatusInProgress}
	attemptRepo.attempts[attempt.ID] = attempt

	return newFakeAnswerRepo(), attemptRepo, questionRepo, optionRepo, attempt
}

func TestSubmitAnswer_Success(t *testing.T) {
	answerRepo, attemptRepo, questionRepo, optionRepo, attempt := setupAnswerFixture()
	svc := NewAnswerService(answerRepo, attemptRepo, questionRepo, optionRepo)

	answer, err := svc.SubmitAnswer(context.Background(), attempt.ID, "q1", "opt-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if answer.SelectedOptionID != "opt-1" {
		t.Errorf("expected 'opt-1', got %q", answer.SelectedOptionID)
	}
}

func TestSubmitAnswer_QuestionFromAnotherQuiz_Rejected(t *testing.T) {
	answerRepo, attemptRepo, questionRepo, optionRepo, attempt := setupAnswerFixture()
	svc := NewAnswerService(answerRepo, attemptRepo, questionRepo, optionRepo)

	_, err := svc.SubmitAnswer(context.Background(), attempt.ID, "q-other-quiz", "opt-1")

	if err == nil {
		t.Fatal("expected an error for a question belonging to another quiz, got nil")
	}
}

func TestSubmitAnswer_OptionFromAnotherQuestion_Rejected(t *testing.T) {
	answerRepo, attemptRepo, questionRepo, optionRepo, attempt := setupAnswerFixture()
	svc := NewAnswerService(answerRepo, attemptRepo, questionRepo, optionRepo)

	_, err := svc.SubmitAnswer(context.Background(), attempt.ID, "q1", "opt-other-question")
	if err == nil {
		t.Fatal("expected ans error for an option belonging to another question, got nil")
	}
}

func TestUpdateAnswer_Success(t *testing.T){
	answerRepo, attemptRepo, questionRepo, optionRepo, attempt := setupAnswerFixture()
	svc := NewAnswerService(answerRepo, attemptRepo, questionRepo, optionRepo)

	if _, err := svc.SubmitAnswer(context.Background(), attempt.ID, "q1", "opt-2"); err!=nil{
		t.Fatalf("setup: unexpected error: %v", err)
	}

	updated, err := svc.UpdateAnswer(context.Background(),attempt.ID, "q1", "opt-1")

	if err!=nil{
		t.Fatalf("expected no error, got %v", err)
	}
	if updated.SelectedOptionID != "opt-1"{
		t.Errorf("expected 'opt-1, got %q", updated.SelectedOptionID)
	}

	all, err := answerRepo.FindByAttemptID(context.Background(), attempt.ID)
	if len(all) != 1{
		t.Errorf("expected exactly one answer for this attempt+question, got %d", len(all))
	}
}
