package helpers

import (
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/fsm"
	"github.com/go-telegram/ui/keyboard/reply"

	"github.com/insigmo/asisa_clinic_finder/internal/db"
	"github.com/insigmo/asisa_clinic_finder/internal/local_models"
)

func SendMessage(params *local_models.BaseParams, text string) error {
	msg, err := params.TgBot.SendMessage(params.Ctx, &bot.SendMessageParams{
		ChatID:    params.UserID,
		Text:      text,
		ParseMode: models.ParseModeMarkdown,
	})
	if err != nil {
		return fmt.Errorf("send message '%s' failed: %v", text, err)
	}

	params.Log.Debug(msg.Text)

	return nil
}

func SendMessageWithReplyMarkup(params *local_models.BaseParams, text string, keyboard *reply.ReplyKeyboard) error {
	msg, err := params.TgBot.SendMessage(params.Ctx, &bot.SendMessageParams{
		ChatID:      params.UserID,
		Text:        text,
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: keyboard,
	})
	if err != nil {
		return fmt.Errorf("send message '%s' failed: %v", text, err)
	}

	params.Log.Debug(msg.Text)

	return nil
}

func GetDbManager(params *local_models.BaseParams) *db.Manager {
	dbManager, ok := params.Ctx.Value(local_models.DBManagerKey).(*db.Manager)
	if !ok {
		params.Log.Error("db manager not found in context")
		return nil
	}
	return dbManager
}

func FetchUser(params *local_models.BaseParams, dbManager *db.Manager) *db.User {
	user, err := dbManager.GetUser(params.Ctx, params.UserID)
	userInfo := params.Update.Message.From

	if err != nil || user == nil {
		user = &db.User{
			ID:           userInfo.ID,
			Username:     userInfo.Username,
			Name:         userInfo.FirstName,
			Lastname:     userInfo.LastName,
			IsBot:        userInfo.IsBot,
			City:         params.Update.Message.Text,
			LanguageCode: userInfo.LanguageCode,
			State:        "",
		}
		err = dbManager.InsertOrUpdateUser(params.Ctx, user)
		if err != nil {
			params.Log.Error(err.Error())
			return nil
		}
	}
	return user
}

func SetUserState(params *local_models.BaseParams, state fsm.StateID) {
	dbManager := GetDbManager(params)
	user := FetchUser(params, dbManager)
	if user == nil {
		params.Log.Error("User not found")
		return
	}

	user.State = string(state)
	if err := dbManager.InsertOrUpdateUser(params.Ctx, user); err != nil {
		params.Log.Error(err.Error())
	}
}
