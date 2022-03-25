package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main() {
	webhook := os.Getenv("webhook")
	token := os.Getenv("token")
	debug := os.Getenv("debug")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	if debug == "true" {
		bot.Debug = true
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

	commands, err := bot.GetMyCommands()
	if err == nil {
		log.Println(commands)
	}

	link, err := url.Parse(webhook)
	update := bot.ListenForWebhook(link.Path)

	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case u := <-update:
				handleUpdate(u)
			case <-ticker.C:
				// do stuff
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	s := &http.Server{
		Addr: ":9180",
	}
	err = s.ListenAndServe()
	if err != nil {
		log.Printf("started failed %v", err)
	}

}

func handleUpdate(update tgbotapi.Update) {

	if update.Message.IsCommand() {
		HandleCommand(update.Message)
	}
}
