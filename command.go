package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/niko/boox-bot/user"
	"log"
	"os"
)

const (
	Help  = "help"
	Login = "login"
	Code  = "code"
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
	u := user.Get(message.From.ID)
	boox := NewBoox(u)
	switch message.Command() {
	case Help:
		msg := tgbotapi.NewMessage(message.Chat.ID, "1. Use /code your_email command to get the login code. \n "+
			"2. User /login your_code to login. \n"+
			"3. Then you can send or forward your book to me.")
		_, err := bot.Send(msg)
		if err != nil {
			log.Println("send Help command error.", err)
			return
		}
	case Login:
		code := message.CommandArguments()
		u := user.Get(message.From.ID)
		err, t, uid := boox.LoginBoox(u.Email, code)
		msg := tgbotapi.NewMessage(message.Chat.ID, "Now you can send/forward me your books.")

		if err != nil {
			msg = tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("login failed %s", err))
			return
		} else {
			err := user.UpdateToken(message.From.ID, uid, t)
			if err != nil {
				msg = tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("login failed %s", err))
			}

		}
		_, _ = bot.Send(msg)
	case Code:

		err := boox.Send(message.CommandArguments())
		txt := "OK,email sent.Please check your email and send the code back to me."
		if err != nil {
			txt = fmt.Sprintln(err)
		}
		msg := tgbotapi.NewMessage(message.Chat.ID, txt)

		user.Add(message.From.ID, message.CommandArguments())

		_, _ = bot.Send(msg)
	default:

	}
}

//Upload file
//TODO limit the document type,such as only epub,pdf,mobi ...
func Upload(message *tgbotapi.Message) {
	u := user.Get(message.From.ID)
	boox := NewBoox(u)

	if len(boox.token()) == 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Please login first.")
		_, _ = bot.Send(msg)
		return
	}

	file, err := bot.GetFile(tgbotapi.FileConfig{
		FileID: message.Document.FileID,
	})
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Upload failed.")
		_, _ = bot.Send(msg)

		return
	}

	fileUrl := file.Link(bot.Token)

	err, m := boox.Upload(fileUrl, message.Document.FileName)
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Upload with %s", m))

	if err != nil {
		msg = tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Upload failed %s", err))
	}
	_, _ = bot.Send(msg)

}
