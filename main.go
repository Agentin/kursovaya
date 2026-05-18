package main

import (
	"log"
	"net/http"
	"os/exec"
	"phishing-trainer/handlers"
	"phishing-trainer/models"
	"runtime"
	"time"
)

func main() {
	if err := models.InitDB(); err != nil {
		log.Fatal("Database init failed:", err)
	}

	http.HandleFunc("/", handlers.Portal)
	http.HandleFunc("/fishing", handlers.Index)
	http.HandleFunc("/prewarning/vk", handlers.PrewarningVK)
	http.HandleFunc("/prewarning/ok", handlers.PrewarningOK)
	http.HandleFunc("/vk_login", handlers.VKLogin)
	http.HandleFunc("/ok_login", handlers.OKLogin)
	http.HandleFunc("/submit", handlers.Submit)
	http.HandleFunc("/training_completed", handlers.TrainingCompleted)
	http.HandleFunc("/stats", handlers.Stats)

	// Новые маршруты портала
	http.HandleFunc("/portal", handlers.Portal)
	http.HandleFunc("/test", handlers.Test)
	http.HandleFunc("/test_details", handlers.TestDetails)
	http.HandleFunc("/email_trainer", handlers.EmailTrainer)

	go func() {
		time.Sleep(500 * time.Millisecond)
		openBrowser("http://localhost:8080/portal")
	}()

	log.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func openBrowser(url string) {
	var cmd string
	var args []string
	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default:
		cmd = "xdg-open"
		args = []string{url}
	}
	if err := exec.Command(cmd, args...).Start(); err != nil {
		log.Printf("Не удалось открыть браузер: %v", err)
	}
}
