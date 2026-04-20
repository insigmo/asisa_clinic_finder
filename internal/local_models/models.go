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
	log, _ := ctx.Value("logger").(*zap.Logger)
	if log == nil {
		log = zap.NewNop()
	}

	var userID int64
	if update != nil && update.Message != nil {
		userID = update.Message.Chat.ID
	}

	return &BaseParams{
		Ctx:    ctx,
		TgBot:  tgBot,
		Update: update,
		Log:    log,
		UserID: userID,
	}
}

type Clinic struct {
	Name        string
	Direction   string
	Address     string
	Distance    int
	PhoneNumber string
	PostalCode  int
}

type Province struct {
	ProvinceID int `json:"province"`
}

type Place struct {
	PlaceId string `json:"placeId"`
}
