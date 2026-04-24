package handlers

import (
	"context"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/insigmo/asisa_clinic_finder/internal/constants"
	"github.com/insigmo/asisa_clinic_finder/internal/db"
	"github.com/insigmo/asisa_clinic_finder/internal/helpers"
	"github.com/insigmo/asisa_clinic_finder/internal/i18n"
	"github.com/insigmo/asisa_clinic_finder/internal/keyboards"
	"github.com/insigmo/asisa_clinic_finder/internal/model"
)

func Default(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	params := model.NewBaseParams(ctx, tgBot, update)
	dbManager := helpers.GetDbManager(params)
	user := helpers.FetchUser(params, dbManager)
	if user == nil {
		params.Log.Error("User not found")
		return
	}

	if action, ok := keyboards.ResolveMainMenuAction(params); ok {
		switch action {
		case keyboards.ActionFindClinic:
			RequestClinicDirection(ctx, tgBot, update)
			return
		case keyboards.ActionChangeCity:
			ChangeCity(ctx, tgBot, update)
			return
		case keyboards.ActionChangeLanguage:
			ChangeLanguage(ctx, tgBot, update)
			return
		}
	}

	switch user.State {
	case string(constants.StateChangeCity):
		handleCityInput(ctx, tgBot, update)
		return
	case string(constants.StateStart):
		handleStart(params, user)
		return
	}
	FindClinic(ctx, tgBot, update)
}

func handleStart(params *model.BaseParams, user *db.User) {
	localizer := i18n.New(user.LanguageCode)

	if err := helpers.SendMessageWithReplyMarkup(
		params,
		localizer.SaveUserMessage(),
		keyboards.BuildMainMenu(params.TgBot, user.LanguageCode),
	); err != nil {
		params.Log.Error(err.Error())
		return
	}
}

func handleCityInput(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	params := model.NewBaseParams(ctx, tgBot, update)
	dbManager := helpers.GetDbManager(params)
	user := helpers.FetchUser(params, dbManager)
	if user == nil {
		params.Log.Error("User not found")
		return
	}

	city := strings.TrimSpace(update.Message.Text)
	if city == "" {
		params.Log.Error("Unknown city")
		return
	}

	postalCodes, err := dbManager.FindCity(ctx, city)
	if err != nil {
		params.Log.Error(err.Error())
		return
	}

	localizer := i18n.New(user.LanguageCode)
	if len(postalCodes) == 0 {
		if err = helpers.SendMessage(params, localizer.WrongCityMessage()); err != nil {
			params.Log.Error(err.Error())
		}

		helpers.SetUserState(params, constants.StateIdle)
		return
	}

	user.City = city
	user.State = string(constants.StateIdle)

	if err = dbManager.InsertOrUpdateUser(ctx, user); err != nil {
		params.Log.Error(err.Error())
		return
	}

	if err = helpers.SendMessage(params, localizer.SaveUserMessage()); err != nil {
		params.Log.Error(err.Error())
		return
	}

	if err = keyboards.SendMainMenu(ctx, tgBot, update); err != nil {
		params.Log.Error(err.Error())
	}
}
