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
	"github.com/insigmo/asisa_clinic_finder/internal/logger"
	"github.com/insigmo/asisa_clinic_finder/internal/middleware"
	"github.com/insigmo/asisa_clinic_finder/internal/model"
)

func main() {
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		_, _ = fmt.Fprintln(os.Stderr, "BOT_TOKEN environment variable is not set")
		os.Exit(1)
	}

	log, err := logger.New()
	if err != nil {
		panic(fmt.Errorf("configure logger: %w", err))
	}
	defer func() { _ = log.Sync() }()

	log.Info("starting tg bot")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	ctx = context.WithValue(ctx, model.LoggerKey, log)

	dbManager, err := db.New(ctx)
	if err != nil {
		panic(fmt.Errorf("connect to db: %w", err))
	}
	defer func() { _ = dbManager.Close() }()

	opts := []bot.Option{
		bot.WithMiddlewares(middleware.DB(dbManager)),

		bot.WithMessageTextHandler("/start", bot.MatchTypeExact, handlers.Start),

		bot.WithMessageTextHandler(keyboards.RuCommandFindClinic, bot.MatchTypeExact, handlers.RequestClinicDirection),
		bot.WithMessageTextHandler(keyboards.EsCommandFindClinic, bot.MatchTypeExact, handlers.RequestClinicDirection),
		bot.WithMessageTextHandler(keyboards.EnCommandFindClinic, bot.MatchTypeExact, handlers.RequestClinicDirection),

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
		panic(fmt.Errorf("create telegram bot: %w", err))
	}

	tgBot.Start(ctx)
}
