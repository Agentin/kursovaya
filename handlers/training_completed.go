package handlers

import (
	"net/http"
	"phishing-trainer/models"
	"phishing-trainer/storage"
)

func TrainingCompleted(w http.ResponseWriter, r *http.Request) {
	visitID := r.URL.Query().Get("visit_id")
	site := r.URL.Query().Get("site")
	if visitID == "" || site == "" {
		http.NotFound(w, r)
		return
	}

	// Загружаем все данные сайта
	results, _ := storage.LoadSimulationResults(site)
	avStatsList, _ := storage.LoadAVWarningStats(site)

	var userData *models.SimulationResult
	for _, res := range results {
		if res.VisitID == visitID {
			userData = &res
			break
		}
	}
	if userData == nil {
		http.NotFound(w, r)
		return
	}

	stats := getStats(results)
	avStats := getAVStats(avStatsList)
	totalSuccessful := stats.SuccessfulAvoidance + avStats.LeftSuccessfully

	phishingPercent := 0.0
	successPercent := 0.0
	legitimatePercent := 0.0
	fakePercent := 0.0

	if stats.TotalVisits > 0 {
		phishingPercent = float64(stats.PhishingAttempts) / float64(stats.TotalVisits) * 100
		successPercent = float64(totalSuccessful) / float64(stats.TotalVisits) * 100
	}
	if stats.SubmittedForms > 0 {
		legitimatePercent = float64(stats.LegitimateCredentialsUsed) / float64(stats.SubmittedForms) * 100
		fakePercent = float64(stats.FakeCredentialsUsed) / float64(stats.SubmittedForms) * 100
	}

	data := models.TrainingData{
		UserData:          userData,
		Stats:             stats,
		AVStats:           avStats,
		TotalSuccessful:   totalSuccessful,
		PhishingPercent:   phishingPercent,
		SuccessPercent:    successPercent,
		LegitimatePercent: legitimatePercent,
		FakePercent:       fakePercent,
	}

	renderTemplate(w, "training_completed", data)
}
