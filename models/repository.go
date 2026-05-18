package models

import (
	"database/sql"
)

// ---- Simulation (VK/OK) ----
func InsertSimulationResult(sessionID, visitID, userIP, userAgent string) error {
	_, err := DB.Exec(
		"INSERT INTO simulation_results (session_id, visit_id, user_ip, user_agent) VALUES ($1, $2, $3, $4)",
		sessionID, visitID, userIP, userAgent,
	)
	return err
}

func UpdateSimulationResultAfterSubmit(visitID string, submittedData string, isLegitimate bool) error {
	_, err := DB.Exec(
		"UPDATE simulation_results SET submitted_data=$1, was_submitted=true, is_legitimate=$2, is_phishing_attempt=true WHERE visit_id=$3",
		submittedData, isLegitimate, visitID,
	)
	return err
}

func GetSimulationResultByVisitID(visitID string) (*SimulationResult, error) {
	row := DB.QueryRow(`
        SELECT id, session_id, visit_id, timestamp, submitted_data, was_submitted, is_legitimate, is_phishing_attempt, user_ip, user_agent
        FROM simulation_results WHERE visit_id=$1`, visitID)
	var r SimulationResult
	var submitted sql.NullString
	err := row.Scan(&r.ID, &r.SessionID, &r.VisitID, &r.Timestamp, &submitted, &r.WasSubmitted, &r.IsLegitimate, &r.IsPhishingAttempt, &r.UserIP, &r.UserAgent)
	if err != nil {
		return nil, err
	}
	if submitted.Valid {
		r.SubmittedData = &submitted.String
	}
	return &r, nil
}

func GetAllSimulationResults() ([]SimulationResult, error) {
	rows, err := DB.Query(`
        SELECT id, session_id, visit_id, timestamp, submitted_data, was_submitted, is_legitimate, is_phishing_attempt, user_ip, user_agent
        FROM simulation_results`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []SimulationResult
	for rows.Next() {
		var r SimulationResult
		var submitted sql.NullString
		if err := rows.Scan(&r.ID, &r.SessionID, &r.VisitID, &r.Timestamp, &submitted, &r.WasSubmitted, &r.IsLegitimate, &r.IsPhishingAttempt, &r.UserIP, &r.UserAgent); err != nil {
			continue
		}
		if submitted.Valid {
			r.SubmittedData = &submitted.String
		}
		results = append(results, r)
	}
	return results, nil
}

// ---- AV Warning Stats ----
func InsertAVWarningStat(sessionID, visitID string) error {
	_, err := DB.Exec("INSERT INTO av_warning_stats (session_id, visit_id, warning_shown) VALUES ($1, $2, true)", sessionID, visitID)
	return err
}

func UpdateAVWarningStat(visitID string, setUserLeft, setUserIgnored bool) error {
	if setUserLeft {
		_, err := DB.Exec("UPDATE av_warning_stats SET user_left = true WHERE visit_id = $1", visitID)
		return err
	} else if setUserIgnored {
		_, err := DB.Exec("UPDATE av_warning_stats SET user_ignored_warning = true WHERE visit_id = $1", visitID)
		return err
	}
	return nil
}

func GetAllAVWarningStats() ([]AVWarningStat, error) {
	rows, err := DB.Query("SELECT id, session_id, visit_id, warning_shown, user_left, user_ignored_warning FROM av_warning_stats")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var stats []AVWarningStat
	for rows.Next() {
		var s AVWarningStat
		rows.Scan(&s.ID, &s.SessionID, &s.VisitID, &s.WarningShown, &s.UserLeft, &s.UserIgnoredWarning)
		stats = append(stats, s)
	}
	return stats, nil
}

// ---- Legitimate credentials ----
func IsLegitimateCredential(username, password string) (bool, error) {
	var id int
	err := DB.QueryRow("SELECT id FROM legitimate_credentials WHERE username=$1 AND password=$2", username, password).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// ---- Test questions ----
func GetQuestionsByLevel(level string) ([]TestQuestion, error) {
	rows, err := DB.Query("SELECT id, question_text, option_a, option_b, option_c, option_d, correct_option FROM test_questions WHERE level=$1 ORDER BY id", level)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var qs []TestQuestion
	for rows.Next() {
		var q TestQuestion
		if err := rows.Scan(&q.ID, &q.QuestionText, &q.OptionA, &q.OptionB, &q.OptionC, &q.OptionD, &q.CorrectOption); err != nil {
			continue
		}
		q.Level = level
		qs = append(qs, q)
	}
	return qs, nil
}

// ---- Test attempts ----
func InsertTestAttempt(sessionID, level string, score float64, total int, passed bool) (int, error) {
	var id int
	err := DB.QueryRow(`
        INSERT INTO test_attempts (session_id, level, score, total_questions, passed, completed_at)
        VALUES ($1, $2, $3, $4, $5, NOW()) RETURNING id`,
		sessionID, level, score, total, passed,
	).Scan(&id)
	return id, err
}

func InsertTestAnswer(attemptID, questionID int, selected string, isCorrect bool) error {
	_, err := DB.Exec(`
        INSERT INTO test_answers (attempt_id, question_id, selected_option, is_correct)
        VALUES ($1, $2, $3, $4)`,
		attemptID, questionID, selected, isCorrect,
	)
	return err
}

func GetTestAttempt(attemptID, sessionID string) (*TestAttempt, error) {
	row := DB.QueryRow(`
        SELECT id, session_id, level, score, total_questions, passed, completed_at
        FROM test_attempts WHERE id=$1 AND session_id=$2`, attemptID, sessionID)
	var a TestAttempt
	err := row.Scan(&a.ID, &a.SessionID, &a.Level, &a.Score, &a.TotalQuestions, &a.Passed, &a.CompletedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func GetAttemptDetails(attemptID int) ([]AttemptDetail, error) {
	rows, err := DB.Query(`
        SELECT q.question_text, q.option_a, q.option_b, q.option_c, q.option_d, q.correct_option, a.selected_option, a.is_correct
        FROM test_answers a
        JOIN test_questions q ON a.question_id = q.id
        WHERE a.attempt_id = $1`, attemptID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var details []AttemptDetail
	for rows.Next() {
		var d AttemptDetail
		if err := rows.Scan(&d.QuestionText, &d.OptionA, &d.OptionB, &d.OptionC, &d.OptionD, &d.CorrectOption, &d.SelectedOption, &d.IsCorrect); err != nil {
			continue
		}
		details = append(details, d)
	}
	return details, nil
}

// Получить последние попытки для уровня
func GetUserAttempts(sessionID, level string, limit int) ([]TestAttempt, error) {
	rows, err := DB.Query(`
        SELECT id, session_id, level, score, total_questions, passed, completed_at
        FROM test_attempts WHERE session_id=$1 AND level=$2 ORDER BY completed_at DESC LIMIT $3`, sessionID, level, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var attempts []TestAttempt
	for rows.Next() {
		var a TestAttempt
		rows.Scan(&a.ID, &a.SessionID, &a.Level, &a.Score, &a.TotalQuestions, &a.Passed, &a.CompletedAt)
		attempts = append(attempts, a)
	}
	return attempts, nil
}

// Общая статистика тестов (для всех)
func GetGlobalTestStats() (totalAttempts, totalPassed int, avgScore float64, err error) {
	row := DB.QueryRow(`
        SELECT COUNT(*), COUNT(CASE WHEN passed THEN 1 END), COALESCE(AVG(score), 0)
        FROM test_attempts`)
	err = row.Scan(&totalAttempts, &totalPassed, &avgScore)
	return
}

// ---- Email simulation ----
func InsertEmailSimulation(sessionID, email, simType string) error {
	_, err := DB.Exec(`
        INSERT INTO email_simulation_stats (session_id, email_provided, simulation_type)
        VALUES ($1, $2, $3)`, sessionID, email, simType)
	return err
}

func UpdateEmailClick(sessionID string) error {
	_, err := DB.Exec(`
        UPDATE email_simulation_stats SET clicked_link = true
        WHERE id = (SELECT id FROM email_simulation_stats WHERE session_id = $1 ORDER BY created_at DESC LIMIT 1)`, sessionID)
	return err
}
