package service

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"math/rand"
	"strconv"
	"strings"
)

func SizeMembers(bot *tg.BotAPI, update tg.Update) {
	chatID := update.Message.Chat.ID
	roleMsg := tg.NewMessage(chatID, "Количество игроков:")
	roleMsg.ReplyMarkup = ChooseMembersBtn()
	_, errRole := bot.Send(roleMsg)
	if errRole != nil {
		log.Printf("service.SizeMembers - can't show role: %v", errRole)
	}
}

func ChooseMembers(bot *tg.BotAPI, update tg.Update, store *Store, places *Places) {
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID
	data := update.CallbackQuery.Data
	size, _ := strings.CutPrefix(data, "set_")
	mems, err := strconv.ParseInt(size, 10, 64)
	if err != nil {
		log.Printf("service.ChooseMembers - can't parse: %v", err)
	}
	members := int(mems)
	seed := rand.New(rand.NewSource(chatID))
	round := Round{
		SpyID:   seed.Intn(members),
		Place:   places.Get(seed.Intn(len(places.places))),
		Members: members,
		Roles:   []int{},
		seed:    seed,
	}
	store.Set(chatID, round)
	deleteMsg := tg.NewDeleteMessage(chatID, messageID)
	_, errDlt := bot.Send(deleteMsg)
	if errDlt != nil {
		log.Printf("service.ChooseMembers - can't delete: %v", errDlt)
	}

	roleMsg := tg.NewMessage(chatID, "Игра началась. Разберите роли:")
	roleMsg.ReplyMarkup = ShowRolesBtn(round)
	_, errRole := bot.Send(roleMsg)
	if errRole != nil {
		log.Printf("service.ChooseMembers - can't show role: %v", errRole)
	}
}

func ShowRoles(bot *tg.BotAPI, update tg.Update, store *Store) {
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID
	deleteMsg := tg.NewDeleteMessage(chatID, messageID)
	_, errDlt := bot.Send(deleteMsg)
	if errDlt != nil {
		log.Printf("service.SizeMembers - can't show role: %v", errDlt)
	}
	data := update.CallbackQuery.Data
	size, _ := strings.CutPrefix(data, "role_")
	roles, err := strconv.ParseInt(size, 10, 64)
	if err != nil {
		log.Printf("service.ShowRoles - can't parse: %v", err)
	}
	roleID := int(roles)
	round, ok := store.Get(chatID)
	if !ok {
		roundMsg := tg.NewMessage(chatID, "Упс.. Раунт не найден. Начните заново")
		_, errRole := bot.Send(roundMsg)
		if errRole != nil {
			log.Printf("service.SizeMembers - can't show round: %v", errRole)
		}
		return
	}
	round.Roles = append(round.Roles, roleID)
	store.Set(chatID, round)
	place := round.Place.name
	if round.SpyID == roleID {
		place = "Шпион"
	}
	placeMsg := tg.NewMessage(chatID, place)
	placeMsg.ReplyMarkup = HideBtn()
	_, errPlace := bot.Send(placeMsg)
	if errPlace != nil {
		log.Printf("service.SizeMembers - can't show pl: %v", errPlace)
	}
}

func HideMessage(bot *tg.BotAPI, update tg.Update, store *Store) {
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID
	deleteMsg := tg.NewDeleteMessage(chatID, messageID)
	_, errDlt := bot.Send(deleteMsg)
	if errDlt != nil {
		log.Printf("service.HideMessage - can't delete role: %v", errDlt)
	}
	round, ok := store.Get(chatID)
	if !ok {
		roundMsg := tg.NewMessage(chatID, "Упс.. Раунт не найден. Начните заново")
		_, errRole := bot.Send(roundMsg)
		if errRole != nil {
			log.Printf("service.HideMessage - can't show round: %v", errRole)
		}
		return
	}
	if len(round.Roles) < round.Members {
		roleMsg := tg.NewMessage(chatID, "Игра началась. Разберите роли:")
		roleMsg.ReplyMarkup = ShowRolesBtn(round)
		_, errRole := bot.Send(roleMsg)
		if errRole != nil {
			log.Printf("can't show game over: %v", errRole)
		}
	}
	if len(round.Roles) == round.Members {
		roleMsg := tg.NewMessage(chatID, "Игра началась")
		roleMsg.ReplyMarkup = StopGameBtn()
		_, errRole := bot.Send(roleMsg)
		if errRole != nil {
			log.Printf("service.SizeMembers - can't start: %v", errRole)
		}
	}
}

func StopGame(bot *tg.BotAPI, update tg.Update, store *Store) {
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID
	deleteMsg := tg.NewDeleteMessage(chatID, messageID)
	_, errDlt := bot.Send(deleteMsg)
	if errDlt != nil {
		log.Printf("service.StopGame - stop: %v", errDlt)
	}
	round, ok := store.Get(chatID)
	if !ok {
		roundMsg := tg.NewMessage(chatID, "Упс.. Раунт не найден. Начните заново")
		_, errRole := bot.Send(roundMsg)
		if errRole != nil {
			log.Printf("service.StopGame - can't round: %v", errRole)
		}
		return
	}
	roleMsg := tg.NewMessage(
		chatID,
		fmt.Sprintf("Игра завершина. Шпион: Игрок №%d, место: %s",
			round.SpyID+1, round.Place.name),
	)
	roleMsg.ReplyMarkup = NewGameBtn()
	_, errRole := bot.Send(roleMsg)
	if errRole != nil {
		log.Printf("service.StopGame. Error: %v", errRole)
	}
}

func NewGame(bot *tg.BotAPI, update tg.Update, store *Store, places *Places) {
	var chatID int64
	callback := update.CallbackQuery
	if callback == nil {
		chatID = update.Message.Chat.ID
	}
	if callback != nil {
		chatID = update.CallbackQuery.Message.Chat.ID
	}
	round, ok := store.Get(chatID)
	if !ok {
		membersMsg := tg.NewMessage(chatID, "Вначале задайте количество игроков")
		_, errRole := bot.Send(membersMsg)
		if errRole != nil {
			log.Printf("service.NewGame - can't show role: %v", errRole)
		}
		return
	}
	newRound := Round{
		SpyID:   round.seed.Intn(round.Members),
		Place:   places.Get(round.seed.Intn(len(places.places))),
		Members: round.Members,
		Roles:   []int{},
		seed:    round.seed,
	}
	store.Set(chatID, newRound)
	roleMsg := tg.NewMessage(chatID, "Игра началась. Разберите роли:")
	roleMsg.ReplyMarkup = ShowRolesBtn(newRound)
	_, errRole := bot.Send(roleMsg)
	if errRole != nil {
		log.Printf("service.NewGame - can't show role: %v", errRole)
	}
}

func Start(bot *tg.BotAPI, update tg.Update) {
	msg := tg.NewMessage(update.Message.Chat.ID,
		`🎲 *Игра "Шпион"*

				Это весёлая игра для компании, реализованная через Telegram-бота.
				
				👥 *Как играть:*
				1. Укажите количество игроков.
				2. Бот случайным образом выбирает одно место (например, "школа").
				3. Всем игрокам, кроме одного, отправляется это место в личные сообщения. Один игрок получает слово "Шпион" — он не знает, где все находятся!
				4. Игроки по очереди задают друг другу вопросы о месте. Например: _"Часто ли ты здесь бываешь?"_
				5. Цель — найти шпиона, не выдав при этом само место.
				6. После круга вопросов начинается голосование за шпиона.
				7. Если шпион угадан правильно — побеждают остальные. Если шпион не раскрыт — он побеждает. Но у него есть шанс выиграть даже в случае разоблачения, если он угадает, что за место было загадано.
				
				🕵️ Побеждает внимательность, логика и умение блефовать!`)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = tg.NewReplyKeyboard(
		tg.NewKeyboardButtonRow(
			tg.NewKeyboardButton("🔄 Начать"),
			tg.NewKeyboardButton("💬 Участники"),
			tg.NewKeyboardButton("📝 Описание"),
		),
	)
	keyboard := msg.ReplyMarkup.(tg.ReplyKeyboardMarkup)
	keyboard.ResizeKeyboard = true
	msg.ReplyMarkup = keyboard
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("service.Start - can't show role: %v", err)
	}
}
