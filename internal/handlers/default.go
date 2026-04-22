package handlers

import (
	"context"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/insigmo/asisa_clinic_finder/internal/db"
	"github.com/insigmo/asisa_clinic_finder/internal/fsm_states"
	"github.com/insigmo/asisa_clinic_finder/internal/helpers"
	"github.com/insigmo/asisa_clinic_finder/internal/keyboards"
	"github.com/insigmo/asisa_clinic_finder/internal/local_models"
	"github.com/insigmo/asisa_clinic_finder/internal/localize_manager"
)

func Default(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	params := local_models.NewBaseParams(ctx, tgBot, update)
	dbManager := helpers.GetDbManager(params)
	user := helpers.FetchUser(params, dbManager)

	text := strings.TrimSpace(update.Message.Text)

	if action, ok := keyboards.ResolveMainMenuAction(user.LanguageCode, text); ok {
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
	case string(fsm_states.StateFindClinic):
		FindClinic(ctx, tgBot, update)
		return

	case string(fsm_states.StateChangeCity):
		handleCityInput(ctx, tgBot, update, user)
		return
	}

	localizer := localize_manager.New(user.LanguageCode)
	if err := helpers.SendMessage(params, localizer.StartMessage()); err != nil {
		params.Log.Error(err.Error())
	}
}

func handleCityInput(ctx context.Context, tgBot *bot.Bot, update *models.Update, user *db.User) {
	params := local_models.NewBaseParams(ctx, tgBot, update)
	dbManager := helpers.GetDbManager(params)

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

	localizer := localize_manager.New(user.LanguageCode)
	if len(postalCodes) == 0 {
		if err = helpers.SendMessage(params, localizer.WrongCityMessage()); err != nil {
			params.Log.Error(err.Error())
		}

		user.State = string(fsm_states.StateIdle)
		if err = dbManager.InsertOrUpdateUser(ctx, user); err != nil {
			params.Log.Error(err.Error())
		}

		return
	}

	user.City = city
	user.State = string(fsm_states.StateIdle)

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
