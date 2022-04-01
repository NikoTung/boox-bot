package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/niko/boox-bot/user"
	"os"
)

const (
	Help  = "help"
	Login = "login"
	Code  = "code"
	Me    = "me"
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
		_, _ = bot.Send(msg)

	case Login:
		email := u.Email
		if len(email) == 0 {
			msg := tgbotapi.NewMessage(message.Chat.ID, "I don't have your email yeah. \nUse /help command to get help.")
			_, _ = bot.Send(msg)
			return
		}

		code := message.CommandArguments()
		if len(code) == 0 {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Please tell me your the code you got from login command.")
			_, _ = bot.Send(msg)
			return
		}
		u := user.Get(message.From.ID)
		err, t, uid := boox.LoginBoox(u.Email, code)
		if err != nil {
			msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Login into boox failed,%s", err))
			_, _ = bot.Send(msg)
			return
		}

		err = u.UpdateToken(uid, t)

		if err != nil {
			msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("login failed %s", err))
			_, _ = bot.Send(msg)
			return
		}

		msg := tgbotapi.NewMessage(message.Chat.ID, "Now you can send/forward me your books.")
		_, _ = bot.Send(msg)
	case Code:
		email := message.CommandArguments()
		if len(email) == 0 {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Please tell me your email.")
			_, _ = bot.Send(msg)
			return
		}
		err := boox.Send(email)
		txt := "OK,email sent.Please check your email and send the code back to me."
		if err != nil {
			txt = fmt.Sprintln(err)
		}
		msg := tgbotapi.NewMessage(message.Chat.ID, txt)

		u.UpdateEmail(email)

		_, _ = bot.Send(msg)
	case Me:
		msg := fmt.Sprintf("Not login yeah.")
		if boox.User.IsLogin() {
			msg = fmt.Sprintf("Email %s \n Boox uid %s \n Expire at %d", u.Email, u.BooxUid, u.Expire)
		}
		_, _ = bot.Send(tgbotapi.NewMessage(message.Chat.ID, msg))
	default:

	}
}

//Upload file
//TODO limit the document type,such as only epub,pdf,mobi ...
func Upload(message *tgbotapi.Message) {
	u := user.Get(message.From.ID)
	boox := NewBoox(u)

	if !u.IsLogin() {
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

	err = boox.Upload(fileUrl, message.Document.FileName)
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Upload with %s", err))

	if err != nil {
		msg = tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Upload failed %s", err))
	}
	_, _ = bot.Send(msg)

}
