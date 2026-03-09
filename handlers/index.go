package handlers

import (
	"fmt"
	"net/http"
	"phishing-trainer/models"
	"phishing-trainer/storage"
)

func Index(w http.ResponseWriter, r *http.Request) {
	// Загружаем данные из файлов
	results1, _ := storage.LoadSimulationResults("vk")
	avStats1, _ := storage.LoadAVWarningStats("vk")
	results2, _ := storage.LoadSimulationResults("ok")
	avStats2, _ := storage.LoadAVWarningStats("ok")

	stats1 := getStats(results1)
	stats2 := getStats(results2)
	av1 := getAVStats(avStats1)
	av2 := getAVStats(avStats2)

	successful1 := stats1.SuccessfulAvoidance + av1.LeftSuccessfully
	successful2 := stats2.SuccessfulAvoidance + av2.LeftSuccessfully

	totalVisits := stats1.TotalVisits + stats2.TotalVisits
	totalPhishing := stats1.PhishingAttempts + stats2.PhishingAttempts
	totalAvoidance := successful1 + successful2
	totalLegitimate := stats1.LegitimateCredentialsUsed + stats2.LegitimateCredentialsUsed

	successRate := 0.0
	if totalVisits > 0 {
		successRate = float64(totalAvoidance) / float64(totalVisits) * 100
	}
	progressWidth := fmt.Sprintf("%.1f%%", successRate)
	data := struct {
		Stats1          models.Stats
		AVStats1        models.AVStats
		Stats2          models.Stats
		AVStats2        models.AVStats
		Successful1     int
		Successful2     int
		TotalVisits     int
		TotalPhishing   int
		TotalAvoidance  int
		TotalLegitimate int
		SuccessRate     float64
		ProgressWidth   string
	}{
		Stats1:          stats1,
		AVStats1:        av1,
		Stats2:          stats2,
		AVStats2:        av2,
		Successful1:     successful1,
		Successful2:     successful2,
		TotalVisits:     totalVisits,
		TotalPhishing:   totalPhishing,
		TotalAvoidance:  totalAvoidance,
		TotalLegitimate: totalLegitimate,
		SuccessRate:     successRate,
		ProgressWidth:   progressWidth,
	}

	renderTemplate(w, "index", data)
}
