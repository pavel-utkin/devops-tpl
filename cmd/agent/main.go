// Клиент для сбора runtime метрик и отправки на сервер
package main

import (
	"devops-tpl/internal/agent"
	"devops-tpl/internal/agent/config"
)

// @Title Client-server metrics
// @Description Сервис сбора и хранения метрик
// @Version 1.0

// @Contact.email pavel@utkin-pro.ru

// @contact.name Efim
// @contact.url https://t.me/utkin_pawka
// @contact.email pavel@utkin-pro.ru

// @Tag.name Update
// @Tag.description "Группа запросов обновления метрик"

// @Tag.name Value
// @Tag.description "Группа запросов получения значений метрик"

// @Tag.name Static
// @Tag.description "Группа эндпоинтов со статикой"

func main() {
	config := config.LoadConfig()
	app := agent.NewHTTPClient(config)
	app.Run()
}
