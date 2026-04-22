package handlers

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/insigmo/asisa_clinic_finder/internal/keyboards"
	"github.com/insigmo/asisa_clinic_finder/internal/local_models"
)

func ChangeLanguage(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	params := local_models.NewBaseParams(ctx, tgBot, update)

	if err := keyboards.SendLanguageMenu(ctx, tgBot, update); err != nil {
		params.Log.Error(err.Error())
	}
}

func ChangeLanguageSelection(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	keyboards.HandleLanguageSelection(ctx, tgBot, update)
}
