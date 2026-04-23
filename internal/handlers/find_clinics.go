package handlers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/insigmo/asisa_clinic_finder/internal/fsm_states"
	"github.com/insigmo/asisa_clinic_finder/internal/helpers"
	"github.com/insigmo/asisa_clinic_finder/internal/local_models"
	"github.com/insigmo/asisa_clinic_finder/internal/localize_manager"
	"github.com/insigmo/asisa_clinic_finder/internal/services/clinic"
)

const findClinicTimeout = 15 * time.Second

func RequestClinicDirection(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	params := local_models.NewBaseParams(ctx, tgBot, update)
	dbManager := helpers.GetDbManager(params)
	user := helpers.FetchUser(params, dbManager)

	localizator := localize_manager.New(user.LanguageCode)

	user.State = string(fsm_states.StateFindClinic)

	if err := dbManager.InsertOrUpdateUser(ctx, user); err != nil {
		params.Log.Error(err.Error())
		return
	}

	if err := helpers.SendMessage(params, localizator.AskDirectionMessage()); err != nil {
		params.Log.Error(err.Error())
	}
}

func FindClinic(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	params := local_models.NewBaseParams(ctx, tgBot, update)
	dbManager := helpers.GetDbManager(params)
	user := helpers.FetchUser(params, dbManager)

	localizator := localize_manager.New(user.LanguageCode)

	// TODO добавить поиск поликлиник по Левенштейну
	direction := strings.TrimSpace(update.Message.Text)
	if direction == "" {
		return
	}

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

		user.State = string(fsm_states.StateIdle)
		if uerr := dbManager.InsertOrUpdateUser(ctx, user); uerr != nil {
			params.Log.Error(uerr.Error())
		}

		return
	}

	city, err := validator.TakeCity(opCtx, update.Message.Chat.ID)
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

		user.State = string(fsm_states.StateIdle)
		if uerr := dbManager.InsertOrUpdateUser(ctx, user); uerr != nil {
			params.Log.Error(uerr.Error())
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
		return
	}

	user.State = string(fsm_states.StateIdle)
	if err = dbManager.InsertOrUpdateUser(ctx, user); err != nil {
		params.Log.Error(err.Error())
	}
}
