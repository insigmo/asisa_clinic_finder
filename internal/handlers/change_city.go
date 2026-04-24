package handlers

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/insigmo/asisa_clinic_finder/internal/constants"
	"github.com/insigmo/asisa_clinic_finder/internal/helpers"
	"github.com/insigmo/asisa_clinic_finder/internal/i18n"
	"github.com/insigmo/asisa_clinic_finder/internal/model"
)

func ChangeCity(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	params := model.NewBaseParams(ctx, tgBot, update)
	dbManager := helpers.GetDbManager(params)
	user := helpers.FetchUser(params, dbManager)
	if user == nil {
		params.Log.Error("User not found")
		return
	}
	userInfo := update.Message.From
	localizator := i18n.New(user.LanguageCode)

	if user.LanguageCode == "" {
		user.LanguageCode = userInfo.LanguageCode
	}

	helpers.SetUserState(params, constants.StateChangeCity)

	if err := helpers.SendMessage(params, localizator.ChangeCityMessage()); err != nil {
		params.Log.Error(err.Error())
	}
}
