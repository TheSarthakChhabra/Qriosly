package service

import (
	"context"
	"errors"
	"quiz-backend/model"
	"time"
)

type fakeQuizRepo struct{ quizzes map[string]model.Quiz }

func newFakeQuizRepo() *fakeQuizRepo { return &fakeQuizRepo{quizzes: map[string]model.Quiz{}} }

func (f *fakeQuizRepo) Save(ctx context.Context, q model.Quiz) error {
	f.quizzes[q.ID] = q
	return nil
}

func (f *fakeQuizRepo) FindAll(ctx context.Context) ([]model.Quiz, error) {
	var out []model.Quiz
	for _, q := range f.quizzes {
		out = append(out, q)
	}
	return out, nil
}

func (f *fakeQuizRepo) FindByID(ctx context.Context, id string) (model.Quiz, error) {
	q, ok := f.quizzes[id]
	if !ok {
		return model.Quiz{}, errors.New("quiz not found")
	}
	return q, nil
}

type fakeQuestionRepo struct{ questions map[string]model.Question }

func newFakeQuestionRepo() *fakeQuestionRepo {
	return &fakeQuestionRepo{questions: map[string]model.Question{}}
}

func (f *fakeQuestionRepo) Save(ctx context.Context, q model.Question) error {
	f.questions[q.ID] = q
	return nil
}

func (f *fakeQuestionRepo) FindByQuizID(ctx context.Context, quizID string) ([]model.Question, error) {
	var out []model.Question
	for _, q := range f.questions {
		if q.QuizID == quizID {
			out = append(out, q)
		}
	}
	return out, nil
}

func (f *fakeQuestionRepo) FindByID(ctx context.Context, id string) (model.Question, error) {
	q, ok := f.questions[id]
	if !ok {
		return model.Question{}, errors.New("question not found")
	}
	return q, nil
}

type fakeOptionRepo struct{ options map[string]model.Option }

func newFakeOptionRepo() *fakeOptionRepo { return &fakeOptionRepo{options: map[string]model.Option{}} }

func (f *fakeOptionRepo) Save(ctx context.Context, o model.Option) error {
	f.options[o.ID] = o
	return nil
}

func (f *fakeOptionRepo) FindByQuestionID(ctx context.Context, questionID string) ([]model.Option, error) {
	var out []model.Option
	for _, o := range f.options {
		if o.QuestionID == questionID {
			out = append(out, o)
		}
	}
	return out, nil
}

func (f *fakeOptionRepo) FindByID(ctx context.Context, id string) (model.Option, error) {
	o, ok := f.options[id]
	if !ok {
		return model.Option{}, errors.New("option not found")
	}
	return o, nil
}

type fakeAttemptRepo struct{ attempts map[string]model.Attempt }

func newFakeAttemptRepo() *fakeAttemptRepo {
	return &fakeAttemptRepo{attempts: map[string]model.Attempt{}}
}

func (f *fakeAttemptRepo) Save(ctx context.Context, a model.Attempt) error {
	f.attempts[a.ID] = a
	return nil
}

func (f *fakeAttemptRepo) FindByID(ctx context.Context, id string) (model.Attempt, error) {
	a, ok := f.attempts[id]
	if !ok {
		return model.Attempt{}, errors.New("attempt not found")
	}
	return a, nil
}

func (f *fakeAttemptRepo) MarkSubmitted(ctx context.Context, id, status string, submittedAt time.Time, score int) error {
	a, ok := f.attempts[id]
	if !ok {
		return errors.New("attempt not found")
	}
	a.Status = status
	sa, sc := submittedAt, score
	a.SubmittedAt, a.Score = &sa, &sc
	f.attempts[id] = a
	return nil
}

func (f *fakeAttemptRepo) FindByUserID(ctx context.Context, userID string) ([]model.Attempt, error) {
	var out []model.Attempt
	for _, a := range f.attempts {
		if a.UserID == userID {
			out = append(out, a)
		}
	}
	return out, nil
}

func (f *fakeAttemptRepo) FindByQuizID(ctx context.Context, quizID string) ([]model.Attempt, error) {
	var out []model.Attempt
	for _, a := range f.attempts {
		if a.QuizID == quizID {
			out = append(out, a)
		}
	}
	return out, nil
}

func (f *fakeAttemptRepo) FindSummariesByQuizID(ctx context.Context, quizID string) ([]model.AttemptSummary, error) {
	return nil, nil
}

type fakeAnswerRepo struct{ answers map[string]model.Answer }

func newFakeAnswerRepo() *fakeAnswerRepo { return &fakeAnswerRepo{answers: map[string]model.Answer{}} }

func (f *fakeAnswerRepo) Save(ctx context.Context, a model.Answer) error {
	f.answers[a.ID] = a
	return nil
}

func (f *fakeAnswerRepo) FindAttemptAndQuestion(ctx context.Context, attemptID, questionID string) (model.Answer, error) {
	for _, a := range f.answers {
		if a.AttemptID == attemptID {
			return a, nil
		}
	}
	return model.Answer{}, errors.New("answer not found")
}

func (f *fakeAnswerRepo) UpdateSelectedOption(ctx context.Context, answerID, selectedOptionID string) error {
	a, ok := f.answers[answerID]
	if !ok {
		return errors.New("answer not found")
	}
	a.SelectedOptionID = selectedOptionID
	f.answers[answerID] = a
	return nil
}

func (f *fakeAnswerRepo) FindByAttemptID(ctx context.Context, attemptID string) ([]model.Answer, error) {
	var out []model.Answer
	for _, a := range f.answers {
		if a.AttemptID == attemptID {
			out = append(out, a)
		}
	}
	return out, nil
}

func (f *fakeAnswerRepo) FindByAttemptAndQuestion(ctx context.Context, attemptID, questionID string) (model.Answer, error) {
	for _, a := range f.answers {
		if a.AttemptID == attemptID && a.QuestionID == questionID {
			return a, nil
		}
	}
	return model.Answer{}, errors.New("answer not found")
}