package models

import "time"

// --- Simulation (VK/OK тренажёры) ---
type SimulationResult struct {
	ID                int
	SessionID         string
	VisitID           string
	Timestamp         time.Time
	SubmittedData     *string
	WasSubmitted      bool
	IsLegitimate      bool
	IsPhishingAttempt bool
	UserIP            string
	UserAgent         string
}

type AVWarningStat struct {
	ID                 int
	SessionID          string
	VisitID            string
	WarningShown       bool
	UserLeft           bool
	UserIgnoredWarning bool
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

// --- Тесты и портал ---
type TestQuestion struct {
	ID            int
	Level         string
	QuestionText  string
	OptionA       string
	OptionB       string
	OptionC       string
	OptionD       string
	CorrectOption string
}

type TestAttempt struct {
	ID             int
	SessionID      string
	Level          string
	Score          float64
	TotalQuestions int
	Passed         bool
	CompletedAt    time.Time
}

type TestAnswer struct {
	ID             int
	AttemptID      int
	QuestionID     int
	SelectedOption string
	IsCorrect      bool
}

type AttemptDetail struct {
	QuestionText   string
	OptionA        string
	OptionB        string
	OptionC        string
	OptionD        string
	CorrectOption  string
	SelectedOption string
	IsCorrect      bool
}

// --- Тренажёр писем ---
type EmailSimulationStats struct {
	ID             int
	SessionID      string
	EmailProvided  *string
	SimulationType string
	ClickedLink    bool
	CreatedAt      time.Time
}
