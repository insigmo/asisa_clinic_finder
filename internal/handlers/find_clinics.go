package handlers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/insigmo/asisa_clinic_finder/internal/db"
	"github.com/insigmo/asisa_clinic_finder/internal/helpers"
	"github.com/insigmo/asisa_clinic_finder/internal/local_models"
	"github.com/insigmo/asisa_clinic_finder/internal/localize_manager"
	"github.com/insigmo/asisa_clinic_finder/internal/services/clinic"
)

const findClinicTimeout = 15 * time.Second

func FindClinic(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	params := local_models.NewBaseParams(ctx, tgBot, update)
	localizator := localize_manager.New(update.Message.From.LanguageCode)

	dbManager, ok := ctx.Value(local_models.DBManagerKey).(*db.Manager)
	if !ok {
		params.Log.Error("dbManager is not set to context")
		return
	}

	// Всё после команды — это направление (может быть многословным).
	raw := strings.TrimSpace(update.Message.Text)
	parts := strings.SplitN(raw, " ", 2)
	if len(parts) < 2 || strings.TrimSpace(parts[1]) == "" {
		return
	}
	direction := strings.TrimSpace(parts[1])

	opCtx, cancel := context.WithTimeout(ctx, findClinicTimeout)
	defer cancel()

	validator := clinic.NewDirectionValidator(dbManager)
	if err := validator.ValidateDirection(opCtx, direction); err != nil {
		similar, ferr := validator.FindSimilarDirections(opCtx, direction)
		if ferr != nil {
			params.Log.Error(ferr.Error())
			return
		}
		msg := localizator.WrongDirection()
		if len(similar) > 0 {
			msg += localizator.Perhaps() + strings.Join(similar, ", ")
		}
		if serr := helpers.SendMessage(params, msg); serr != nil {
			params.Log.Error(serr.Error())
		}
		return
	}

	userID := update.Message.Chat.ID
	city, err := validator.TakeCity(opCtx, userID)
	if err != nil {
		params.Log.Error(fmt.Sprintf("Failed when tried to get user: %v", err))
		return
	}

	postalCodes, err := dbManager.FindCity(opCtx, city)
	if err != nil {
		params.Log.Error(err.Error())
		return
	}
	if len(postalCodes) == 0 {
		if serr := helpers.SendMessage(params, localizator.WrongCityMessage()); serr != nil {
			params.Log.Error(serr.Error())
		}
		return
	}

	postalCode := postalCodes[len(postalCodes)/2]
	result, err := clinic.Search(opCtx, city, direction, postalCode)
	if err != nil {
		params.Log.Error(err.Error())
		return
	}
	if err = helpers.SendMessage(params, result); err != nil {
		params.Log.Error(err.Error())
	}
}
