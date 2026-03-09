package handlers

import (
	"encoding/json"
	"net/http"
	"phishing-trainer/models"
	"phishing-trainer/storage"
)

func Submit(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	r.ParseForm()
	visitID := r.FormValue("visit_id")
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Определяем сайт по префиксу visit_id (можно по-другому, но для примера)
	// В реальности лучше передавать site через скрытое поле или из сессии.
	// Упростим: будем искать в обоих хранилищах.
	var site string
	results1, _ := storage.LoadSimulationResults("vk")
	found := false
	for _, res := range results1 {
		if res.VisitID == visitID {
			site = "vk"
			found = true
			break
		}
	}
	if !found {
		results2, _ := storage.LoadSimulationResults("ok")
		for _, res := range results2 {
			if res.VisitID == visitID {
				site = "ok"
				found = true
				break
			}
		}
	}
	if !found {
		http.Error(w, "Invalid visit ID", http.StatusBadRequest)
		return
	}

	// Проверяем легитимность (хардкод для демо)
	isLegitimate := (username == "ivan@example.com" && password == "12345") ||
		(username == "sasha" && password == "password")

	submittedJSON, _ := json.Marshal(map[string]string{"username": username, "password": password})
	dataStr := string(submittedJSON)

	// Обновляем запись
	storage.UpdateSimulationResult(site, visitID, func(r *models.SimulationResult) {
		r.SubmittedData = &dataStr
		r.WasSubmitted = true
		r.IsLegitimate = isLegitimate
		r.IsPhishingAttempt = true
	})

	http.Redirect(w, r, "/training_completed?visit_id="+visitID+"&site="+site, http.StatusSeeOther)
}
