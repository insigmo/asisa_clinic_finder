package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"

	"github.com/insigmo/asisa_clinic_finder/internal/db"
	"github.com/insigmo/asisa_clinic_finder/internal/handlers"
	"github.com/insigmo/asisa_clinic_finder/internal/keyboards"
	"github.com/insigmo/asisa_clinic_finder/internal/local_models"
	"github.com/insigmo/asisa_clinic_finder/internal/logger"
	"github.com/insigmo/asisa_clinic_finder/internal/middlewares"
)

const BotToken = "8583496833:AAGZTPUXNqlcopLRZ2KqJ2myW0-muVE7WSo"

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		token = BotToken
	}

	log, err := logger.ConfigureLogger()
	if err != nil {
		panic(err)
	}
	defer func() { _ = log.Sync() }()

	log.Info("Starting tg bot")

	ctx = context.WithValue(ctx, local_models.LoggerKey, log)

	dbManager, err := db.New(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to connect to db: %w", err))
	}
	defer func() { _ = dbManager.Close() }()

	opts := []bot.Option{
		bot.WithMiddlewares(middlewares.DBMiddleware(dbManager)),
		bot.WithMessageTextHandler("/start", bot.MatchTypeExact, handlers.Start),

		bot.WithMessageTextHandler(keyboards.RuCommandFindClinic, bot.MatchTypePrefix, handlers.RequestClinicDirection),
		bot.WithMessageTextHandler(keyboards.EsCommandFindClinic, bot.MatchTypePrefix, handlers.RequestClinicDirection),
		bot.WithMessageTextHandler(keyboards.EnCommandFindClinic, bot.MatchTypePrefix, handlers.RequestClinicDirection),

		bot.WithMessageTextHandler(keyboards.RuCommandChangeCity, bot.MatchTypeExact, handlers.ChangeCity),
		bot.WithMessageTextHandler(keyboards.EsCommandChangeCity, bot.MatchTypeExact, handlers.ChangeCity),
		bot.WithMessageTextHandler(keyboards.EnCommandChangeCity, bot.MatchTypeExact, handlers.ChangeCity),

		bot.WithMessageTextHandler(keyboards.RuCommandChangeLanguage, bot.MatchTypeExact, handlers.ChangeLanguage),
		bot.WithMessageTextHandler(keyboards.EsCommandChangeLanguage, bot.MatchTypeExact, handlers.ChangeLanguage),
		bot.WithMessageTextHandler(keyboards.EnCommandChangeLanguage, bot.MatchTypeExact, handlers.ChangeLanguage),

		bot.WithDefaultHandler(handlers.Default),
	}

	tgBot, err := bot.New(token, opts...)
	if err != nil {
		panic(fmt.Errorf("failed to create telegram bot: %w", err))
	}

	tgBot.Start(ctx)
}
