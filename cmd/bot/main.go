package main

import (
	"context"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/clients/abstract"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/clients/telegram"
	config2 "gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/config"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/model/callbacks"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/model/messages"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/repository"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/service"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
		<-exit
		cancel()
	}()

	config, err := config2.New()
	if err != nil {
		log.Fatal("config init failed", err)
	}

	tgClient, err := telegram.New(config)
	if err != nil {
		log.Fatal("telegram client init failed", err)
	}

	abstractClient := abstract.NewCurrencyClient(config)

	transactionRepo, err := repository.NewTransactionRepository()
	if err != nil {
		log.Fatal("transaction repository init failed", err)
	}

	userCurrencyRepo, err := repository.NewUserCurrencyRepository(config)
	if err != nil {
		log.Fatal("user currency repository init failed", err)
	}

	exchangeRatesService, err := service.NewExchangeRatesService(ctx, abstractClient)
	if err != nil {
		log.Fatal("exchange rates service init failed", err)
	}

	financeCalculatorService := service.NewFinanceCalculatorService(transactionRepo, exchangeRatesService)

	msgModel := messages.New(tgClient, userCurrencyRepo)
	callbackModel := callbacks.New(tgClient, transactionRepo, userCurrencyRepo, exchangeRatesService, financeCalculatorService)

	tgClient.ListenUpdates(msgModel, callbackModel)
}
