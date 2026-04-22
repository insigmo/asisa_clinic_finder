package helpers

import (
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
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

	if err != nil || user == nil {
		params.Log.Error("user not found")
		return &db.User{}
	}
	return user
}
