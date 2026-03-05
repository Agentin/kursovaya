package handlers

import (
	"html/template"
	"net/http"
	"phishing-trainer/models"
)

func Submit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Парсим данные формы
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Ошибка обработки формы", http.StatusBadRequest)
		return
	}

	submission := models.Submission{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
		IP:       r.RemoteAddr,
	}

	// В реальном тренажёре можно сохранять в лог или БД, но для примера просто покажем
	renderTemplate(w, "result", submission)
}

// Вспомогательная функция для рендеринга шаблонов
func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	templates := template.Must(template.ParseFiles(
		"templates/base.html",
		"templates/"+tmpl+".html",
	))
	err := templates.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}