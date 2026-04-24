package helpers

import (
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/reply"

	"github.com/insigmo/asisa_clinic_finder/internal/constants"
	"github.com/insigmo/asisa_clinic_finder/internal/db"
	"github.com/insigmo/asisa_clinic_finder/internal/model"
)

// SendMessage отправляет текстовое сообщение пользователю (Markdown).
func SendMessage(params *model.BaseParams, text string) error {
	msg, err := params.TgBot.SendMessage(params.Ctx, &bot.SendMessageParams{
		ChatID:    params.UserID,
		Text:      text,
		ParseMode: models.ParseModeMarkdown,
	})
	if err != nil {
		return fmt.Errorf("send message: %w", err)
	}
	params.Log.Debug(msg.Text)
	return nil
}

// SendMessageWithReplyMarkup отправляет сообщение с reply-клавиатурой.
func SendMessageWithReplyMarkup(params *model.BaseParams, text string, keyboard *reply.ReplyKeyboard) error {
	msg, err := params.TgBot.SendMessage(params.Ctx, &bot.SendMessageParams{
		ChatID:      params.UserID,
		Text:        text,
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: keyboard,
	})
	if err != nil {
		return fmt.Errorf("send message with markup: %w", err)
	}
	params.Log.Debug(msg.Text)
	return nil
}

// GetDbManager извлекает *db.Manager из контекста.
func GetDbManager(params *model.BaseParams) *db.Manager {
	dbManager, ok := params.Ctx.Value(model.DBManagerKey).(*db.Manager)
	if !ok {
		params.Log.Error("db manager not found in context")
		return nil
	}
	return dbManager
}

// FetchUser возвращает пользователя из БД, создавая запись если её нет.
func FetchUser(params *model.BaseParams, dbManager *db.Manager) *db.User {
	user, err := dbManager.GetUser(params.Ctx, params.UserID)
	if err != nil {
		params.Log.Error("get user: " + err.Error())
		return nil
	}
	if user != nil {
		return user
	}

	info := params.Update.Message.From
	newUser := &db.User{
		ID:           info.ID,
		Username:     info.Username,
		Name:         info.FirstName,
		Lastname:     info.LastName,
		IsBot:        info.IsBot,
		LanguageCode: info.LanguageCode,
	}
	if err = dbManager.InsertOrUpdateUser(params.Ctx, newUser); err != nil {
		params.Log.Error("insert user: " + err.Error())
		return nil
	}
	return newUser
}

// SetUserState сохраняет новое состояние FSM пользователя в БД.
func SetUserState(params *model.BaseParams, state constants.StateID) {
	dbManager := GetDbManager(params)
	if dbManager == nil {
		return
	}
	user := FetchUser(params, dbManager)
	if user == nil {
		params.Log.Error("set user state: user not found")
		return
	}
	user.State = string(state)
	if err := dbManager.InsertOrUpdateUser(params.Ctx, user); err != nil {
		params.Log.Error("set user state: " + err.Error())
	}
}
