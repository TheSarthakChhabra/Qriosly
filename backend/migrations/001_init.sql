CREATE TABLE quizzes (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    duration_minutes INTEGER NOT NULL DEFAULT 10,
    created_by UUID,
    CONSTRAINT title_not_empty CHECK (title <> ''),
    CONSTRAINT duration_positive CHECK (duration_minutes > 0)
);

CREATE TABLE users (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'student',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT name_not_empty CHECK (name <> ''),
    CONSTRAINT email_not_empty CHECK (email <> ''),
    CONSTRAINT valid_role CHECK (role IN ('student', 'teacher', 'admin'))
);

ALTER TABLE quizzes ADD CONSTRAINT quizzes_created_by_fkey FOREIGN KEY (created_by) REFERENCES users(id);

CREATE TABLE questions (
    id UUID PRIMARY KEY,
    quiz_id UUID NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE,
    text TEXT NOT NULL,
    CONSTRAINT question_text_not_empty CHECK (text <> '')
);

CREATE TABLE options (
    id UUID PRIMARY KEY,
    question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    text TEXT NOT NULL,
    is_correct BOOLEAN NOT NULL DEFAULT false,
    CONSTRAINT option_text_not_empty CHECK (text <> '')
);

CREATE TABLE attempts (
    id UUID PRIMARY KEY,
    quiz_id UUID NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    status TEXT NOT NULL DEFAULT 'in_progress',
    submitted_at TIMESTAMPTZ,
    score INTEGER,
    CONSTRAINT valid_status CHECK (status IN ('in_progress', 'submitted'))
);

CREATE TABLE answers (
    id UUID PRIMARY KEY,
    attempt_id UUID NOT NULL REFERENCES attempts(id) ON DELETE CASCADE,
    question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    selected_option_id UUID NOT NULL REFERENCES options(id) ON DELETE CASCADE,
    CONSTRAINT unique_attempt_question UNIQUE (attempt_id, question_id)
);