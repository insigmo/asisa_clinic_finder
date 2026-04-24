package handlers

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/insigmo/asisa_clinic_finder/internal/keyboards"
	"github.com/insigmo/asisa_clinic_finder/internal/model"
)

func ChangeLanguage(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	params := model.NewBaseParams(ctx, tgBot, update)

	if err := keyboards.SendLanguageMenu(ctx, tgBot, update); err != nil {
		params.Log.Error(err.Error())
	}
}
