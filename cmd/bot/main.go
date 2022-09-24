package main

import (
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/clients/tg"
	config2 "gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/config"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/model/messages"
	"log"
)

func main() {
	config, err := config2.New()
	if err != nil {
		log.Fatal("config init failed", err)
	}

	tgClient, err := tg.New(config)
	if err != nil {
		log.Fatal("tg client init failed", err)
	}

	msgModel := messages.New(tgClient)

	tgClient.ListenUpdates(msgModel)
}
