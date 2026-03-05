package handlers

import (
	"net/http"
)

func PhishingPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "phishing", nil)
}