package service

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func ChooseMembersBtn() tg.InlineKeyboardMarkup {
	buttons := []string{"3", "4", "5", "6", "7", "8", "9", "10", "11"}
	var rows [][]tg.InlineKeyboardButton
	for i := 0; i < len(buttons); i++ {
		btn := tg.NewInlineKeyboardButtonData(
			buttons[i], "set_"+buttons[i],
		)
		rows = append(rows, []tg.InlineKeyboardButton{btn})
	}
	return tg.NewInlineKeyboardMarkup(rows...)

}

func HideBtn() tg.InlineKeyboardMarkup {
	hide := tg.NewInlineKeyboardButtonData(
		"Скрыть", "hide",
	)
	keyboard := tg.NewInlineKeyboardMarkup(
		[]tg.InlineKeyboardButton{hide},
	)
	return keyboard
}

func contains(slice []int, num int) bool {
	for _, v := range slice {
		if v == num {
			return true
		}
	}
	return false
}

func StopGameBtn() tg.InlineKeyboardMarkup {
	keyboard := tg.NewInlineKeyboardMarkup()
	keyboard.InlineKeyboard = [][]tg.InlineKeyboardButton{
		{
			tg.NewInlineKeyboardButtonData("Завершить", "stop"),
		},
	}
	return keyboard
}

func NewGameBtn() tg.InlineKeyboardMarkup {
	keyboard := tg.NewInlineKeyboardMarkup()
	keyboard.InlineKeyboard = [][]tg.InlineKeyboardButton{
		{
			tg.NewInlineKeyboardButtonData("Новая игра", "new"),
		},
	}
	return keyboard
}

func ShowRolesBtn(round Round) tg.InlineKeyboardMarkup {
	keyboard := tg.NewInlineKeyboardMarkup()
	var rows [][]tg.InlineKeyboardButton
	for i := 0; i < round.Members; i++ {
		if !contains(round.Roles, i) {
			btn := tg.NewInlineKeyboardButtonData(
				fmt.Sprintf("Игрок №%d", i+1),
				fmt.Sprintf("role_%d", i),
			)
			rows = append(rows, []tg.InlineKeyboardButton{btn})
		}
	}
	keyboard.InlineKeyboard = rows
	return keyboard
}
