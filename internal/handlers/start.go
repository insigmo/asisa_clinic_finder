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

func Start(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	params := model.NewBaseParams(ctx, tgBot, update)
	dbManager := helpers.GetDbManager(params)
	user := helpers.FetchUser(params, dbManager)
	if user == nil {
		params.Log.Error("User not found")
		return
	}
	languageCode := user.LanguageCode
	localizator := i18n.New(languageCode)

	helpers.SetUserState(params, constants.StateChangeCity)

	if err := helpers.SendMessage(params, localizator.StartMessage()); err != nil {
		params.Log.Error(err.Error())
		return
	}
}
