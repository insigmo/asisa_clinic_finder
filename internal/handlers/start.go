package handlers

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/insigmo/asisa_clinic_finder/internal/helpers"
	"github.com/insigmo/asisa_clinic_finder/internal/keyboards"
	"github.com/insigmo/asisa_clinic_finder/internal/local_models"
	"github.com/insigmo/asisa_clinic_finder/internal/localize_manager"
)

func Start(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	params := local_models.NewBaseParams(ctx, tgBot, update)
	languageCode := update.Message.From.LanguageCode
	localizator := localize_manager.New(languageCode)

	if err := helpers.SendMessage(params, localizator.StartMessage()); err != nil {
		params.Log.Error(err.Error())
		return
	}

	if err := keyboards.SendMainMenu(ctx, tgBot, update); err != nil {
		params.Log.Error(err.Error())
	}
}
