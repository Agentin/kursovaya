package models

import "time"

type SimulationResult struct {
	ID                int       `json:"id"`
	SessionID         string    `json:"session_id"`
	VisitID           string    `json:"visit_id"`
	Timestamp         time.Time `json:"timestamp"`
	SubmittedData     *string   `json:"submitted_data"`
	WasSubmitted      bool      `json:"was_submitted"`
	IsLegitimate      bool      `json:"is_legitimate"`
	IsPhishingAttempt bool      `json:"is_phishing_attempt"`
	UserIP            string    `json:"user_ip"`
	UserAgent         string    `json:"user_agent"`
}

type AVWarningStat struct {
	ID                 int    `json:"id"`
	SessionID          string `json:"session_id"`
	VisitID            string `json:"visit_id"`
	WarningShown       bool   `json:"warning_shown"`
	UserLeft           bool   `json:"user_left"`
	UserIgnoredWarning bool   `json:"user_ignored_warning"`
}

type Stats struct {
	TotalVisits               int
	SubmittedForms            int
	LegitimateCredentialsUsed int
	FakeCredentialsUsed       int
	PhishingAttempts          int
	SuccessfulAvoidance       int
}

type AVStats struct {
	TotalWarnings    int
	LeftSuccessfully int
	IgnoredWarning   int
}

type TrainingData struct {
	UserData          *SimulationResult
	Stats             Stats
	AVStats           AVStats
	TotalSuccessful   int
	PhishingPercent   float64
	SuccessPercent    float64
	LegitimatePercent float64
	FakePercent       float64
}

type StatsPage struct {
	Stats           Stats
	AVStats         AVStats
	TotalSuccessful int
	PhishingPercent float64
	Details         []SimulationResult
}
