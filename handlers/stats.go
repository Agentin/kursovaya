package handlers

import (
	"net/http"
	"phishing-trainer/models"
	"phishing-trainer/storage"
)

func Stats(w http.ResponseWriter, r *http.Request) {
	site := r.URL.Query().Get("site")
	if site == "" {
		site = "vk"
	}

	results, _ := storage.LoadSimulationResults(site)
	avStatsList, _ := storage.LoadAVWarningStats(site)

	stats := getStats(results)
	avStats := getAVStats(avStatsList)
	totalSuccessful := stats.SuccessfulAvoidance + avStats.LeftSuccessfully

	phishingPercent := 0.0
	if stats.TotalVisits > 0 {
		phishingPercent = float64(stats.PhishingAttempts) / float64(stats.TotalVisits) * 100
	}

	// Выбираем только те, где was_submitted = true
	var details []models.SimulationResult
	for _, r := range results {
		if r.WasSubmitted {
			details = append(details, r)
		}
	}

	data := models.StatsPage{
		Stats:           stats,
		AVStats:         avStats,
		TotalSuccessful: totalSuccessful,
		PhishingPercent: phishingPercent,
		Details:         details,
	}

	renderTemplate(w, "stats", data)
}
