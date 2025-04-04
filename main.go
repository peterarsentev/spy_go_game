package main

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	spy "spy_game/service"
	"strings"
)

func main() {
	config, err := spy.ReadConfig("application.properties")
	if err != nil {
		log.Fatal("Error reading config:", err)
	}
	botToken := config.Get("telegram.token.api")
	bot, err := tg.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}

	store := spy.NewStore()
	places := spy.NewPlaces()

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
				spy.ChooseMembers(bot, update, store, places)
			}
			if strings.HasPrefix(data, "role_") {
				spy.ShowRoles(bot, update, store)
			}
			if data == "hide" {
				spy.HideMessage(bot, update, store)
			}
			if data == "stop" {
				spy.StopGame(bot, update, store)
			}
			if data == "new" {
				spy.NewGame(bot, update, store, places)
			}
		}
		if update.Message == nil {
			continue
		}
		if update.Message.Text == "/start" || update.Message.Text == "📝 Описание" {
			spy.Start(bot, update)
		}
		if update.Message.Text == "🔄 Начать" {
			spy.NewGame(bot, update, store, places)
		}
		if update.Message.Text == "💬 Участники" {
			spy.SizeMembers(bot, update)
		}
	}
}
