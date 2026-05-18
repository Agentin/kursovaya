package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"phishing-trainer/models"
)

func EmailTrainer(w http.ResponseWriter, r *http.Request) {
	sess, _ := getSession(r)

	if r.Method == http.MethodPost && r.FormValue("generate_email") != "" {
		email := r.FormValue("email")
		simType := r.FormValue("sim_type")
		if simType == "" {
			simType = "bank"
		}
		models.InsertEmailSimulation(sess.ID, email, simType)

		data := struct {
			ShowEmail bool
			EmailType string
			EmailHTML template.HTML
		}{
			ShowEmail: true,
			EmailType: simType,
			EmailHTML: template.HTML(generateEmailHTML(simType)),
		}
		renderTemplate(w, "email_trainer", data)
		return
	}

	renderTemplate(w, "email_trainer", struct{ ShowEmail bool }{false})
}

// generateEmailHTML возвращает HTML письма в зависимости от типа
// Ссылки ведут на локальные тренажёры
func generateEmailHTML(simType string) string {
	// Определяем целевую ссылку в зависимости от типа письма
	var targetURL string
	if simType == "social" {
		targetURL = "/prewarning/ok"
	} else {
		targetURL = "/prewarning/vk"
	}

	link := fmt.Sprintf(`<a href="%s" style="background: #d32f2f; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px;">Перейти</a>`, targetURL)

	switch simType {
	case "bank":
		return `
            <div style="font-family: Arial, sans-serif; border: 1px solid #ccc; padding: 20px; max-width: 600px;">
                <h2 style="color: #d32f2f;">⚠️ Срочное уведомление от банка</h2>
                <p><strong>Отправитель:</strong> security@bank-holding.com</p>
                <p><strong>Тема:</strong> Ваш счёт будет заблокирован</p>
                <hr>
                <p>Уважаемый клиент,</p>
                <p>Мы обнаружили подозрительную активность в вашем личном кабинете. Для предотвращения блокировки срочно подтвердите свои данные:</p>
                <p style="text-align: center;">` + link + `</p>
                <p style="color: #666; font-size: 12px;">* Это автоматическое сообщение, отвечать на него не нужно.</p>
            </div>`
	case "social":
		return `
            <div style="font-family: Arial, sans-serif; border: 1px solid #ccc; padding: 20px; max-width: 600px;">
                <h2 style="color: #1e88e5;">🔔 Кто-то хочет войти в ваш аккаунт</h2>
                <p><strong>Отправитель:</strong> support@vk-ok.ru</p>
                <p><strong>Тема:</strong> Подозрительный вход</p>
                <hr>
                <p>Здравствуйте!</p>
                <p>Мы зафиксировали попытку входа в ваш аккаунт из нового устройства. Если это были не вы, срочно смените пароль:</p>
                <p style="text-align: center;">` + link + `</p>
            </div>`
	default: // delivery
		return `
            <div style="font-family: Arial, sans-serif; border: 1px solid #ccc; padding: 20px; max-width: 600px;">
                <h2 style="color: #43a047;">📦 Не удалось доставить посылку</h2>
                <p><strong>Отправитель:</strong> tracking@pochta-delivery.ru</p>
                <p><strong>Тема:</strong> Уточните данные для доставки</p>
                <hr>
                <p>Уважаемый клиент,</p>
                <p>К сожалению, мы не можем доставить вашу посылку, так как адрес указан не полностью. Перейдите по ссылке и укажите правильные данные:</p>
                <p style="text-align: center;">` + link + `</p>
            </div>`
	}
}
