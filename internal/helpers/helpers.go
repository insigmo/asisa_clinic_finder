package helpers

import (
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/insigmo/asisa_clinic_finder/internal/local_models"
)

func SendMessage(params *local_models.BaseParams, text string) error {
	msg, err := params.TgBot.SendMessage(params.Ctx, &bot.SendMessageParams{
		ChatID: params.UserID,
		Text:   text,
	})

	if err != nil {
		return fmt.Errorf("send message '%s' failed: %v", text, err)
	}

	params.Log.Debug(msg.Text)

	return nil
}
