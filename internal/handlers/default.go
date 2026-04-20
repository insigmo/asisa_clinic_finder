package handlers

import (
	"context"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/insigmo/asisa_clinic_finder/internal/db"
	"github.com/insigmo/asisa_clinic_finder/internal/helpers"
	"github.com/insigmo/asisa_clinic_finder/internal/local_models"
	"github.com/insigmo/asisa_clinic_finder/internal/localize_manager"
)

func Default(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	text := strings.TrimSpace(update.Message.Text)
	params := local_models.NewBaseParams(ctx, tgBot, update)
	dbManager, ok := ctx.Value(local_models.DBManagerKey).(*db.Manager)

	if !ok {
		params.Log.Error("db manager not found in context")
		return
	}

	userInfo := update.Message.From
	user := &db.User{
		ID:           userInfo.ID,
		Username:     userInfo.Username,
		Name:         userInfo.FirstName,
		Lastname:     userInfo.LastName,
		IsBot:        userInfo.IsBot,
		City:         text,
		LanguageCode: userInfo.LanguageCode,
	}

	if err := dbManager.InsertOrUpdateUser(ctx, user); err != nil {
		params.Log.Error(err.Error())

		return
	}

	localizer := localize_manager.New(userInfo.LanguageCode)
	if err := helpers.SendMessage(params, localizer.SaveUserMessage()); err != nil {
		params.Log.Error(err.Error())
	}
}
