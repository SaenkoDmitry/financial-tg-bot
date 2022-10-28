package main

import (
	"context"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/clients/abstract"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/clients/telegram"
	config2 "gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/config"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/db"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/model/callbacks"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/model/messages"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/repository"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/service"
	"log"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	config, err := config2.New()
	handleError(err, "config init failed")

	// ----- clients -----
	telegramClient, err := telegram.New(config)
	handleError(err, "telegram client init failed")

	abstractClient := abstract.NewCurrencyClient(config)

	// ----- db init -----
	dbPool, err := db.InitPool(config)
	handleError(err, "could set up database")
	defer dbPool.Close()

	// ----- repositories -----
	transactionRepo := repository.NewTransactionRepository(dbPool)
	userRepo := repository.NewUserRepository(dbPool)
	categoryRepo := repository.NewCategoryRepository(dbPool)
	rateRepo := repository.NewRateRepository(dbPool)
	limitationRepo := repository.NewLimitationRepository(dbPool)

	// ----- services -----
	//ratesCache := cache.New(defaultExpiration, cleanupInterval)
	simpleCache := service.NewSimpleCache(ctx, config.RatesCacheDefaultExpiration(), config.RatesCacheCleanupInterval())

	rateService := service.NewCurrencyExchangeService(ctx, abstractClient, simpleCache, rateRepo)

	calcService := service.NewCalculatorService(transactionRepo, rateRepo, rateService)

	// ----- logic -----
	msgModel := messages.New(telegramClient, userRepo, categoryRepo)
	callbackModel := callbacks.New(telegramClient, transactionRepo, userRepo, categoryRepo, limitationRepo,
		rateService, calcService)

	telegramClient.ListenUpdates(ctx, msgModel, callbackModel)
}

func handleError(err error, message string) {
	if err != nil {
		log.Fatal(errors.Wrap(err, message))
	}
}
