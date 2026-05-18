package handlers

import (
	"net/http"
	"phishing-trainer/models"
	"strconv"
)

func TestDetails(w http.ResponseWriter, r *http.Request) {
	attemptIDStr := r.URL.Query().Get("attempt_id")
	if attemptIDStr == "" {
		http.NotFound(w, r)
		return
	}
	attemptID, err := strconv.Atoi(attemptIDStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	sess, _ := getSession(r)

	attempt, err := models.GetTestAttempt(attemptIDStr, sess.ID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	details, err := models.GetAttemptDetails(attemptID)
	if err != nil {
		http.Error(w, "Ошибка загрузки деталей", http.StatusInternalServerError)
		return
	}

	data := struct {
		Attempt *models.TestAttempt
		Details []models.AttemptDetail
	}{
		Attempt: attempt,
		Details: details,
	}
	renderTemplate(w, "test_details", data)
}
