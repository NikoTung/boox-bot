package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

const (
	Help  = "Help"
	Login = "Login"
)

var locationKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("China", "cn"),
		tgbotapi.NewInlineKeyboardButtonData("United States", "us"),
	),
)

var bot *tgbotapi.BotAPI

func init() {
	bot, _ = tgbotapi.NewBotAPI(os.Getenv("token"))

}

func HandleCommand(message *tgbotapi.Message) {
	switch message.Command() {
	case Help:
		msg := tgbotapi.NewMessage(message.Chat.ID, "Use /Login your_email command to Login your boox account first, then you can send or forward your book to this boot.")
		_, err := bot.Send(msg)
		if err != nil {
			log.Println("send Help command error.", err)
			return
		}
	case Login:
		msg := tgbotapi.NewMessage(message.Chat.ID, "Choose your Login location.")

		// If the message was open, add a copy of our numeric keyboard.
		switch message.Text {
		case "open":
			msg.ReplyMarkup = locationKeyboard
		}

		// Send the message.
		if _, err := bot.Send(msg); err != nil {
			panic(err)
		}
	default:

	}
}
