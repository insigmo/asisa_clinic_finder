package handlers

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/insigmo/asisa_clinic_finder/internal/db"
	"github.com/insigmo/asisa_clinic_finder/internal/helpers"
	"github.com/insigmo/asisa_clinic_finder/internal/local_models"
	"github.com/insigmo/asisa_clinic_finder/internal/localize_manager"
)

func ChangeCity(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	params := local_models.NewBaseParams(ctx, tgBot, update)
	localizator := localize_manager.New(update.Message.From.LanguageCode)

	dbManager, ok := ctx.Value(local_models.DBManagerKey).(*db.Manager)
	if !ok {
		params.Log.Error("dbManager is not set to context")
		return
	}

	userID := update.Message.Chat.ID
	user, err := dbManager.GetUser(ctx, userID)
	if err != nil || user == nil {
		params.Log.Error(fmt.Sprintf("Failed when tried to get user: %v", err))
		return
	}

	city := update.Message.Text
	postalCodes, err := dbManager.FindCity(ctx, city)
	if err != nil {
		params.Log.Error(err.Error())
		return
	}

	if len(postalCodes) == 0 {
		if err = helpers.SendMessage(params, localizator.WrongCityMessage()); err != nil {
			params.Log.Error(err.Error())
		}
		return
	}

	user.City = city
	if err = dbManager.InsertOrUpdateUser(ctx, user); err != nil {
		params.Log.Error(err.Error())
		return
	}

	if err = helpers.SendMessage(params, localizator.StartMessage()); err != nil {
		params.Log.Error(err.Error())
	}
}
