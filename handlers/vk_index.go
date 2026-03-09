package handlers

import (
	"net/http"
	"phishing-trainer/models"
	"phishing-trainer/storage"
	"time"
)

func VKLogin(w http.ResponseWriter, r *http.Request) {
	session, _ := getSession(r)
	visitID := generateVisitID()
	session.Values["visit_id"] = visitID
	session.Save(r, w)

	// Создаём запись о визите
	result := models.SimulationResult{
		SessionID:         session.ID,
		VisitID:           visitID,
		Timestamp:         time.Now(),
		WasSubmitted:      false,
		IsPhishingAttempt: false,
		UserIP:            r.RemoteAddr,
		UserAgent:         r.UserAgent(),
	}
	storage.AppendSimulationResult("vk", result)

	data := struct{ VisitID string }{VisitID: visitID}
	renderTemplate(w, "vk_login", data)
}
