package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/url"
	"os"
)

func main() {
	webhook := os.Getenv("webhook")
	token := os.Getenv("token")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	info, err := bot.GetWebhookInfo()
	if err != nil && len(webhook) == 0 {
		log.Panic("webhook not exist.")
	}

	if webhook != info.URL {
		webhookConfig, err := tgbotapi.NewWebhook(webhook)
		if err != nil {
			log.Panic("config webhook failed,", err)
		}

		_, err = bot.Request(webhookConfig)
		if err != nil {
			log.Panic("config webhook to telegram failed,", err)
			return
		}
	}

	link, err := url.Parse(webhook)
	update := bot.ListenForWebhook(link.Path)

	select {
	case x := <-update:
		go handleUpdate(x)
	}
}

func handleUpdate(x tgbotapi.Update) {


}