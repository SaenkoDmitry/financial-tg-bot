package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/mem"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/service"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/tracing"

	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/clients/abstract"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/clients/telegram"
	config2 "gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/config"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/db"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/logger"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/model/callbacks"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/model/messages"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/repository"
	"go.uber.org/zap"
)

func main() {
	tracing.InitTracing()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	config, err := config2.New()
	handleError(err, "config init failed")

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		port := flag.Int("port", 9095, "the port to listen")
		logger.Info("starting http server", zap.Int("port", *port))
		err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
		if err != nil {
			logger.Fatal("error starting http server", zap.Error(err))
		}
	}()

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
	//ratesCache := mem.New(defaultExpiration, cleanupInterval)
	//simpleCache := service.NewSimpleCache(ctx, config.RatesCacheDefaultExpiration(), config.RatesCacheCleanupInterval())
	memcached := mem.NewMemcached(config.CacheHost())

	rateService := service.NewCurrencyExchangeService(ctx, abstractClient, memcached, rateRepo)

	calcService := service.NewCalculatorService(config, transactionRepo, rateRepo, rateService, memcached)

	// ----- logic -----
	msgModel := messages.New(telegramClient, userRepo, categoryRepo)
	callbackModel := callbacks.New(telegramClient, transactionRepo, userRepo, categoryRepo, limitationRepo,
		rateService, calcService, memcached)

	telegramClient.ListenUpdates(ctx, msgModel, callbackModel)
}

func handleError(err error, message string) {
	if err != nil {
		logger.Fatal(message, zap.Error(err))
	}
}
