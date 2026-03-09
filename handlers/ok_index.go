package handlers

import (
	"net/http"
	"phishing-trainer/models"
	"phishing-trainer/storage"
	"time"
)

func OKLogin(w http.ResponseWriter, r *http.Request) {
	session, _ := getSession(r)
	visitID := generateVisitID()
	session.Values["visit_id"] = visitID
	session.Save(r, w)

	// Создаём запись о посещении
	result := models.SimulationResult{
		SessionID:         session.ID,
		VisitID:           visitID,
		Timestamp:         time.Now(),
		WasSubmitted:      false,
		IsPhishingAttempt: false,
		UserIP:            r.RemoteAddr,
		UserAgent:         r.UserAgent(),
	}
	// Сохраняем в хранилище для сайта "ok"
	storage.AppendSimulationResult("ok", result)

	data := struct{ VisitID string }{VisitID: visitID}
	renderTemplate(w, "ok_login", data)
}
