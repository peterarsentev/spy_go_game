package main

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	"spy_game/internal/service"
	"strings"
)

func main() {
	path := "application.properties"
	if len(os.Args) > 0 {
		path = os.Args[1]
	}
	config, err := service.ReadConfig(path)
	if err != nil {
		log.Fatal("Error reading config:", err)
	}
	botToken, ok := config.Get("telegram.token.api")
	if !ok {
		log.Fatal("not found token")
	}
	bot, err := tg.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("not start bot: %s", err)
	}

	store := service.NewStore()
	places := service.NewPlaces()

	fmt.Printf("Go to https://t.me/%s\n", bot.Self.UserName)
	u := tg.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		if update.CallbackQuery != nil {
			data := update.CallbackQuery.Data
			if strings.HasPrefix(data, "set_") {
				service.ChooseMembers(bot, update, store, places)
			}
			if strings.HasPrefix(data, "role_") {
				service.ShowRoles(bot, update, store)
			}
			if data == "hide" {
				service.HideMessage(bot, update, store)
			}
			if data == "stop" {
				service.StopGame(bot, update, store)
			}
			if data == "new" {
				service.NewGame(bot, update, store, places)
			}
		}
		if update.Message == nil {
			continue
		}
		if update.Message.Text == "/start" || update.Message.Text == "📝 Описание" {
			service.Start(bot, update)
		}
		if update.Message.Text == "🔄 Начать" {
			service.NewGame(bot, update, store, places)
		}
		if update.Message.Text == "💬 Участники" {
			service.SizeMembers(bot, update)
		}
	}
}
