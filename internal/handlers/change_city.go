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

var MinLanguageCodeLen = 2

func ChangeCity(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	params := local_models.NewBaseParams(ctx, tgBot, update)
	localizator := localize_manager.New(update.Message.From.LanguageCode)

	dbManager, ok := ctx.Value("dbManager").(*db.Manager)
	if !ok {
		params.Log.Error(fmt.Sprintf("dbManager is not set to context"))

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
		params.Log.Error(fmt.Sprintf("Unknown city"))

		if err = helpers.SendMessage(params, localizator.WrongCityMessage()); err != nil {
			params.Log.Error(err.Error())

			return
		}

		return
	}

	user.City = city

	err = dbManager.InsertOrUpdateUser(ctx, user)
	if err != nil {
		return
	}

	if err = helpers.SendMessage(params, localizator.StartMessage()); err != nil {
		params.Log.Error(err.Error())

		return
	}
}
