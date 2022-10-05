package main

import (
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/clients/tg"
	config2 "gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/config"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/model/callbacks"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/model/messages"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/repository"
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

	transactionRepo, err := repository.NewTransactionRepository()
	if err != nil {
		log.Fatal("transaction repository init failed", err)
	}

	msgModel := messages.New(tgClient)
	callbackModel := callbacks.New(tgClient, transactionRepo)

	tgClient.ListenUpdates(msgModel, callbackModel)
}
