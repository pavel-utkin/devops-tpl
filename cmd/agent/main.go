// Клиент для сбора runtime метрик и отправки на сервер
package main

import (
	"context"
	"devops-tpl/internal/agent"
	"devops-tpl/internal/agent/config"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var buildVersion = "N/A"
var buildDate = "N/A"
var buildCommit = "N/A"

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

	ctx, ctxCancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer ctxCancel()

	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	config := config.LoadConfig()
	app := agent.NewHTTPClient(config)
	app.Run(ctx)
}
