package handlers

import (
	"context"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/insigmo/asisa_clinic_finder/internal/db"
	"github.com/insigmo/asisa_clinic_finder/internal/helpers"
	"github.com/insigmo/asisa_clinic_finder/internal/keyboards"
	"github.com/insigmo/asisa_clinic_finder/internal/local_models"
	"github.com/insigmo/asisa_clinic_finder/internal/localize_manager"
)

func Default(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	text := strings.TrimSpace(update.Message.Text)
	params := local_models.NewBaseParams(ctx, tgBot, update)

	if text == "" {
		return
	}

	dbManager, ok := ctx.Value(local_models.DBManagerKey).(*db.Manager)
	if !ok {
		params.Log.Error("db manager not found in context")
		return
	}

	userInfo := update.Message.From
	userID := update.Message.Chat.ID

	user, err := dbManager.GetUser(ctx, userID)
	if err != nil || user == nil {
		user = &db.User{
			ID:           userInfo.ID,
			Username:     userInfo.Username,
			Name:         userInfo.FirstName,
			Lastname:     userInfo.LastName,
			IsBot:        userInfo.IsBot,
			LanguageCode: userInfo.LanguageCode,
		}
	}

	if user.LanguageCode == "" {
		user.LanguageCode = userInfo.LanguageCode
	}

	if action, ok := keyboards.ResolveMainMenuAction(user.LanguageCode, text); ok {
		switch action {
		case keyboards.ActionFindClinic:
			if err = helpers.SendMessage(params, "Please use /find_clinic <medical direction>"); err != nil {
				params.Log.Error(err.Error())
			}
			return
		case keyboards.ActionChangeCity:
			ChangeCity(ctx, tgBot, update)
			return
		case keyboards.ActionChangeLanguage:
			ChangeLanguage(ctx, tgBot, update)
			return
		}
	}

	if _, ok = keyboards.ResolveLanguage(text); ok {
		ChangeLanguageSelection(ctx, tgBot, update)
		return
	}

	city := text
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
		return
	}

	user.City = city
	if user.LanguageCode == "" {
		user.LanguageCode = userInfo.LanguageCode
	}

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
