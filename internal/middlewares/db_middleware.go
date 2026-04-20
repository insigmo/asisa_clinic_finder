package middlewares

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/insigmo/asisa_clinic_finder/internal/db"
	"github.com/insigmo/asisa_clinic_finder/internal/local_models"
)

func DBMiddleware(dbManager *db.Manager) bot.Middleware {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, update *models.Update) {
			ctx = context.WithValue(ctx, local_models.DBManagerKey, dbManager)
			next(ctx, b, update)
		}
	}
}
