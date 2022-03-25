package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

const (
	help  = "help"
	login = "login"
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
	case help:
		msg := tgbotapi.NewMessage(message.Chat.ID, "Use /login your_email command to login your boox account first, then you can send or forward your book to this boot.")
		_, err := bot.Send(msg)
		if err != nil {
			log.Println("send help command error.", err)
			return
		}
	case login:
		msg := tgbotapi.NewMessage(message.Chat.ID, "Choose your login location.")

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
