package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/insigmo/asisa_clinic_finder/internal/db"
	"github.com/insigmo/asisa_clinic_finder/internal/handlers"
	"github.com/insigmo/asisa_clinic_finder/internal/logger"
	"github.com/insigmo/asisa_clinic_finder/internal/middlewares"
)

const BotToken = "8583496833:AAGZTPUXNqlcopLRZ2KqJ2myW0-muVE7WSo"

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	log, err := logger.ConfigureLogger()
	if err != nil {
		panic(err)
	}

	log.Info("Starting tg bot")

	ctx = context.WithValue(ctx, "logger", log)
	dbManager := db.New(ctx)

	opts := []bot.Option{
		bot.WithMiddlewares(middlewares.DBMiddleware(dbManager)),
		bot.WithMessageTextHandler("/start", bot.MatchTypeExact, handlers.Start),
		bot.WithMessageTextHandler("/find_clinic", bot.MatchTypeExact, handlers.Start),
	}

	tgBot, _ := bot.New(BotToken, opts...)
	tgBot.Start(ctx)
}
