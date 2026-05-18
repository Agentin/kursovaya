package handlers

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"phishing-trainer/models"
	"time"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("super-secret-key"))

func renderTemplate(w http.ResponseWriter, tmplName string, data interface{}) {
	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
	}
	tmpl, err := template.New(tmplName + ".html").Funcs(funcMap).ParseFiles("templates/" + tmplName + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func generateVisitID() string {
	return fmt.Sprintf("visit_%d_%d", time.Now().UnixNano(), rand.Intn(1000))
}

func getSession(r *http.Request) (*sessions.Session, error) {
	return store.Get(r, "phishing-trainer-session")
}

// getStatsFromResults вычисляет статистику из слайса SimulationResult
func getStatsFromResults(results []models.SimulationResult) models.Stats {
	s := models.Stats{}
	for _, r := range results {
		s.TotalVisits++
		if r.WasSubmitted {
			s.SubmittedForms++
			s.PhishingAttempts++
			if r.IsLegitimate {
				s.LegitimateCredentialsUsed++
			} else {
				s.FakeCredentialsUsed++
			}
		} else {
			s.SuccessfulAvoidance++
		}
	}
	return s
}

func getAVStatsFromList(stats []models.AVWarningStat) models.AVStats {
	av := models.AVStats{}
	for _, st := range stats {
		av.TotalWarnings++
		if st.UserLeft {
			av.LeftSuccessfully++
		}
		if st.UserIgnoredWarning {
			av.IgnoredWarning++
		}
	}
	return av
}

// getStats вычисляет статистику из слайса SimulationResult
func getStats(results []models.SimulationResult) models.Stats {
	s := models.Stats{}
	for _, r := range results {
		s.TotalVisits++
		if r.WasSubmitted {
			s.SubmittedForms++
			s.PhishingAttempts++
			if r.IsLegitimate {
				s.LegitimateCredentialsUsed++
			} else {
				s.FakeCredentialsUsed++
			}
		} else {
			s.SuccessfulAvoidance++
		}
	}
	return s
}

// getAVStats вычисляет статистику из слайса AVWarningStat
func getAVStats(stats []models.AVWarningStat) models.AVStats {
	av := models.AVStats{}
	for _, st := range stats {
		av.TotalWarnings++
		if st.UserLeft {
			av.LeftSuccessfully++
		}
		if st.UserIgnoredWarning {
			av.IgnoredWarning++
		}
	}
	return av
}
