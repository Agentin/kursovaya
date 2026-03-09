package handlers

import (
	"net/http"
	"phishing-trainer/models"
	"phishing-trainer/storage"
)

func PrewarningVK(w http.ResponseWriter, r *http.Request) {
	prewarning(w, r, "vk")
}

func PrewarningOK(w http.ResponseWriter, r *http.Request) {
	prewarning(w, r, "ok")
}

func prewarning(w http.ResponseWriter, r *http.Request, site string) {
	session, _ := getSession(r)
	visitID, ok := session.Values["av_warning_visit_id"].(string)
	if !ok {
		visitID = generateVisitID()
		session.Values["av_warning_visit_id"] = visitID
		session.Save(r, w)

		stat := models.AVWarningStat{
			SessionID:    session.ID,
			VisitID:      visitID,
			WarningShown: true,
		}
		storage.AppendAVWarningStat(site, stat)
	}

	if r.Method == "POST" {
		action := r.FormValue("action")
		if action == "proceed" {
			storage.UpdateAVWarningStat(site, visitID, func(s *models.AVWarningStat) {
				s.UserIgnoredWarning = true
			})
			if site == "vk" {
				http.Redirect(w, r, "/vk_login", http.StatusSeeOther)
			} else {
				http.Redirect(w, r, "/ok_login", http.StatusSeeOther)
			}
			return
		} else if action == "leave" {
			storage.UpdateAVWarningStat(site, visitID, func(s *models.AVWarningStat) {
				s.UserLeft = true
			})
			http.Redirect(w, r, "https://www.google.com", http.StatusSeeOther)
			return
		}
	}

	data := struct {
		Site string
		URL  string
	}{
		Site: site,
		URL:  r.Host,
	}
	renderTemplate(w, "prewarning", data)
}
