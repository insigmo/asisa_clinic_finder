package handlers

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/insigmo/asisa_clinic_finder/internal/fsm_states"
	"github.com/insigmo/asisa_clinic_finder/internal/helpers"
	"github.com/insigmo/asisa_clinic_finder/internal/local_models"
	"github.com/insigmo/asisa_clinic_finder/internal/localize_manager"
)

func ChangeCity(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	params := local_models.NewBaseParams(ctx, tgBot, update)
	dbManager := helpers.GetDbManager(params)
	user := helpers.FetchUser(params, dbManager)

	userInfo := update.Message.From
	localizator := localize_manager.New(user.LanguageCode)

	if user.LanguageCode == "" {
		user.LanguageCode = userInfo.LanguageCode
	}

	user.State = string(fsm_states.StateChangeCity)

	if err := dbManager.InsertOrUpdateUser(ctx, user); err != nil {
		params.Log.Error(err.Error())
		return
	}

	if err := helpers.SendMessage(params, localizator.ChangeCityMessage()); err != nil {
		params.Log.Error(err.Error())
	}
}
