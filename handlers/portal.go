package handlers

import (
	"net/http"
	"phishing-trainer/models"
)

func Portal(w http.ResponseWriter, r *http.Request) {
	sess, _ := getSession(r)

	// Прогресс из сессии (для отображения значков)
	progress := map[string]bool{
		"basic":  false,
		"medium": false,
		"expert": false,
	}
	if val, ok := sess.Values["progress_basic"].(bool); ok && val {
		progress["basic"] = true
	}
	if val, ok := sess.Values["progress_medium"].(bool); ok && val {
		progress["medium"] = true
	}
	if val, ok := sess.Values["progress_expert"].(bool); ok && val {
		progress["expert"] = true
	}

	// Получение последних попыток для каждого уровня
	attemptsBasic, _ := models.GetUserAttempts(sess.ID, "basic", 5)
	attemptsMedium, _ := models.GetUserAttempts(sess.ID, "medium", 5)
	attemptsExpert, _ := models.GetUserAttempts(sess.ID, "expert", 5)

	// Общая статистика тестов
	totalAttempts, totalPassed, avgScore, _ := models.GetGlobalTestStats()

	data := struct {
		Progress       map[string]bool
		AttemptsBasic  []models.TestAttempt
		AttemptsMedium []models.TestAttempt
		AttemptsExpert []models.TestAttempt
		TotalAttempts  int
		TotalPassed    int
		AvgScore       float64
	}{
		Progress:       progress,
		AttemptsBasic:  attemptsBasic,
		AttemptsMedium: attemptsMedium,
		AttemptsExpert: attemptsExpert,
		TotalAttempts:  totalAttempts,
		TotalPassed:    totalPassed,
		AvgScore:       avgScore,
	}

	renderTemplate(w, "portal", data)
}
