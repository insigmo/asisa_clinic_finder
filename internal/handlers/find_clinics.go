package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/insigmo/asisa_clinic_finder/internal/db"
	"github.com/insigmo/asisa_clinic_finder/internal/helpers"
	"github.com/insigmo/asisa_clinic_finder/internal/local_models"
	"github.com/insigmo/asisa_clinic_finder/internal/localize_manager"
	"github.com/insigmo/asisa_clinic_finder/internal/services/clinic"
)

func FindClinic(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	params := local_models.NewBaseParams(ctx, tgBot, update)
	localizator := localize_manager.New(update.Message.From.LanguageCode)

	dbManager, ok := ctx.Value("dbManager").(*db.Manager)
	if !ok {
		params.Log.Error(fmt.Sprintf("dbManager is not set to context"))

		return
	}
	directionValidator := clinic.NewDirectionValidator(dbManager)
	text := strings.Split(update.Message.Text, " ")
	if len(text) < 2 {
		return
	}
	direction := text[1]
	err := directionValidator.ValidateDirection(direction)
	if err != nil {
		similarDirections, err := directionValidator.FindSimilarDirections(direction)
		if err != nil {
			params.Log.Error(err.Error())

			return
		}
		msg := localizator.WrongDirection()

		if len(similarDirections) > 0 {
			msg += localizator.Perhaps() + strings.Join(similarDirections, ", ")
		}

		if err = helpers.SendMessage(params, msg); err != nil {
			params.Log.Error(err.Error())

			return
		}
		return
	}

	userID := update.Message.Chat.ID
	city, err := directionValidator.TakeCity(userID)
	if err != nil {
		params.Log.Error(fmt.Sprintf("Failed when tried to get user: %v", err))

		return
	}

	postalCodes, err := dbManager.FindCity(ctx, city)
	if err != nil {
		params.Log.Error(err.Error())

		return
	}

	postalCode := postalCodes[len(postalCodes)/2]
	if err = helpers.SendMessage(params, clinic.Search(city, direction, postalCode)); err != nil {
		params.Log.Error(err.Error())

		return
	}
}
