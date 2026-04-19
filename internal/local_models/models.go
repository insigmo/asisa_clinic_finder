package local_models

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

type BaseParams struct {
	Ctx    context.Context
	TgBot  *bot.Bot
	Update *models.Update
	Log    *zap.Logger
	UserID int64
}

func NewBaseParams(ctx context.Context, tgBot *bot.Bot, update *models.Update) *BaseParams {
	log := ctx.Value("logger").(*zap.Logger)

	return &BaseParams{
		Ctx:    ctx,
		TgBot:  tgBot,
		Update: update,
		Log:    log,
		UserID: update.Message.Chat.ID,
	}
}
