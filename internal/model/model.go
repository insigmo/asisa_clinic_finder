package model

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

// ctxKey — приватный тип для ключей контекста, предотвращает коллизии.
type ctxKey int

const (
	LoggerKey    ctxKey = iota
	DBManagerKey ctxKey = iota
)

// BaseParams — общий набор зависимостей, передаваемых в хендлеры и сервисы.
type BaseParams struct {
	Ctx    context.Context
	TgBot  *bot.Bot
	Update *models.Update
	Log    *zap.Logger
	UserID int64
}

// NewBaseParams собирает BaseParams из аргументов хендлера.
func NewBaseParams(ctx context.Context, tgBot *bot.Bot, update *models.Update) *BaseParams {
	log, _ := ctx.Value(LoggerKey).(*zap.Logger)
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

// Clinic представляет клинику из ответа ASISA.
type Clinic struct {
	Name        string
	Direction   string
	Address     string
	Distance    int
	PhoneNumber string
	PostalCode  int
}

// Province — ответ эндпоинта gmaps-service.
type Province struct {
	ProvinceID int `json:"province"`
}

// Place — элемент ответа autocomplete-addresses.
type Place struct {
	PlaceID string `json:"placeId"`
}
