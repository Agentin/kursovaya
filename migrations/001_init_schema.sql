-- Таблицы для тренажёров ВК и ОК
CREATE TABLE IF NOT EXISTS simulation_results (
    id SERIAL PRIMARY KEY,
    session_id VARCHAR(255),
    visit_id VARCHAR(255) UNIQUE,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    submitted_data TEXT,
    was_submitted BOOLEAN DEFAULT FALSE,
    is_legitimate BOOLEAN DEFAULT FALSE,
    is_phishing_attempt BOOLEAN DEFAULT FALSE,
    user_ip VARCHAR(45),
    user_agent TEXT
);

CREATE TABLE IF NOT EXISTS av_warning_stats (
    id SERIAL PRIMARY KEY,
    session_id VARCHAR(255),
    visit_id VARCHAR(255),
    warning_shown BOOLEAN DEFAULT FALSE,
    user_left BOOLEAN DEFAULT FALSE,
    user_ignored_warning BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS legitimate_credentials (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL
);

-- Таблицы для портала, тестов и тренажёра писем
CREATE TABLE IF NOT EXISTS test_questions (
    id SERIAL PRIMARY KEY,
    level VARCHAR(20) NOT NULL CHECK (level IN ('basic', 'medium', 'expert')),
    question_text TEXT NOT NULL,
    option_a TEXT NOT NULL,
    option_b TEXT NOT NULL,
    option_c TEXT NOT NULL,
    option_d TEXT NOT NULL,
    correct_option CHAR(1) NOT NULL CHECK (correct_option IN ('A', 'B', 'C', 'D'))
);

CREATE TABLE IF NOT EXISTS test_attempts (
    id SERIAL PRIMARY KEY,
    session_id VARCHAR(255) NOT NULL,
    level VARCHAR(20) NOT NULL,
    score FLOAT NOT NULL,
    total_questions INTEGER NOT NULL,
    passed BOOLEAN NOT NULL,
    completed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS test_answers (
    id SERIAL PRIMARY KEY,
    attempt_id INTEGER REFERENCES test_attempts(id) ON DELETE CASCADE,
    question_id INTEGER REFERENCES test_questions(id),
    selected_option CHAR(1),
    is_correct BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS email_simulation_stats (
    id SERIAL PRIMARY KEY,
    session_id VARCHAR(255) NOT NULL,
    email_provided VARCHAR(255),
    simulation_type VARCHAR(50),
    clicked_link BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);