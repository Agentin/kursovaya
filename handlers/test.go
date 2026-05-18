package handlers

import (
	"net/http"
	"phishing-trainer/models"
	"strconv"
)

func Test(w http.ResponseWriter, r *http.Request) {
	level := r.URL.Query().Get("level")
	if level == "" {
		level = "basic"
	}
	if level != "basic" && level != "medium" && level != "expert" {
		http.Error(w, "Invalid level", http.StatusBadRequest)
		return
	}

	sess, _ := getSession(r)

	// Если POST – обрабатываем ответы
	if r.Method == http.MethodPost {
		r.ParseForm()
		questions, _ := models.GetQuestionsByLevel(level)
		total := len(questions)
		correct := 0

		answers := make(map[int]string)
		for qidStr, selected := range r.Form {
			if qidStr == "level" {
				continue
			}
			qid, _ := strconv.Atoi(qidStr)
			answers[qid] = selected[0]
		}

		results := make([]struct {
			Question  models.TestQuestion
			Selected  string
			IsCorrect bool
		}, 0, total)

		for _, q := range questions {
			selected := answers[q.ID]
			isCorrect := (selected == q.CorrectOption)
			if isCorrect {
				correct++
			}
			results = append(results, struct {
				Question  models.TestQuestion
				Selected  string
				IsCorrect bool
			}{q, selected, isCorrect})
		}

		scorePercent := float64(correct) / float64(total) * 100
		passed := scorePercent >= 80.0

		// Сохраняем попытку
		attemptID, err := models.InsertTestAttempt(sess.ID, level, scorePercent, total, passed)
		if err == nil {
			for _, res := range results {
				models.InsertTestAnswer(attemptID, res.Question.ID, res.Selected, res.IsCorrect)
			}
			if passed {
				sess.Values["progress_"+level] = true
				sess.Save(r, w)
			}
		}

		// Перенаправляем на страницу деталей попытки
		http.Redirect(w, r, "/test_details?attempt_id="+strconv.Itoa(attemptID), http.StatusSeeOther)
		return
	}

	// GET – показываем форму теста
	questions, err := models.GetQuestionsByLevel(level)
	if err != nil {
		http.Error(w, "Ошибка загрузки вопросов", http.StatusInternalServerError)
		return
	}

	data := struct {
		Level     string
		Questions []models.TestQuestion
	}{
		Level:     level,
		Questions: questions,
	}
	renderTemplate(w, "test", data)
}
