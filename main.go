package main

import (
	"log"
	"net/http"
	"os/exec"
	"phishing-trainer/handlers"
	"runtime"
	"time"
)

func main() {
	// Статические файлы
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Маршруты
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/phishing", handlers.PhishingPage)
	http.HandleFunc("/submit", handlers.Submit)

	// Открываем браузер через 0.5 секунды (даём серверу время запуститься)
	go func() {
		time.Sleep(500 * time.Millisecond)
		openBrowser("http://localhost:8080")
	}()

	log.Println("🚀 Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// openBrowser открывает URL в браузере по умолчанию (кросс-платформенно)
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
	default: // linux, freebsd и др.
		cmd = "xdg-open"
		args = []string{url}
	}

	if err := exec.Command(cmd, args...).Start(); err != nil {
		log.Printf("❌ Не удалось открыть браузер: %v", err)
	}
}