package handlers

import (
	"context"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/insigmo/asisa_clinic_finder/internal/constants"
	"github.com/insigmo/asisa_clinic_finder/internal/db"
	"github.com/insigmo/asisa_clinic_finder/internal/helpers"
	"github.com/insigmo/asisa_clinic_finder/internal/i18n"
	"github.com/insigmo/asisa_clinic_finder/internal/keyboards"
	"github.com/insigmo/asisa_clinic_finder/internal/model"
	"github.com/insigmo/asisa_clinic_finder/internal/services/clinic"
)

const findClinicTimeout = 15 * time.Second

func RequestClinicDirection(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	params := model.NewBaseParams(ctx, tgBot, update)
	dbManager := helpers.GetDbManager(params)
	user := helpers.FetchUser(params, dbManager)
	if user == nil {
		params.Log.Error("User not found")
		return
	}
	localizator := i18n.New(user.LanguageCode)

	helpers.SetUserState(params, constants.StateIdle)

	if err := helpers.SendMessage(params, localizator.AskDirectionMessage()); err != nil {
		params.Log.Error(err.Error())
	}
}

func FindClinic(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	params := model.NewBaseParams(ctx, tgBot, update)
	dbManager := helpers.GetDbManager(params)
	user := helpers.FetchUser(params, dbManager)
	if user == nil {
		params.Log.Error("User not found")
		return
	}

	direction := strings.TrimSpace(update.Message.Text)
	direction, err := validateDirection(params, dbManager, user, direction)
	if err != nil {
		params.Log.Error(err.Error())
		return
	}

	result, err := clinic.Search(params, dbManager, user, direction)
	if err != nil {
		params.Log.Error(err.Error())
		return
	}

	if err = helpers.SendMessageWithReplyMarkup(params, result, keyboards.BuildMainMenu(params.TgBot, user.LanguageCode)); err != nil {
		params.Log.Error(err.Error())
		return
	}
	helpers.SetUserState(params, constants.StateIdle)
}

func validateDirection(params *model.BaseParams, dbManager *db.Manager, user *db.User, direction string) (string, error) {
	localizator := i18n.New(user.LanguageCode)

	opCtx, cancel := context.WithTimeout(params.Ctx, findClinicTimeout)
	defer cancel()

	validator := clinic.NewDirectionValidator(dbManager)
	validName, err := dbManager.FindMedicalDirection(params.Ctx, direction)

	if err != nil {
		similar, ferr := validator.FindSimilarDirections(opCtx, validName)
		if ferr != nil {
			params.Log.Error(ferr.Error())
			return "", ferr
		}

		msg := localizator.WrongDirection()
		if len(similar) > 0 {
			msg += localizator.Perhaps() + strings.Join(similar, ", ")
		}

		if serr := helpers.SendMessage(params, msg); serr != nil {
			params.Log.Error(serr.Error())
			return "", serr
		}

		helpers.SetUserState(params, constants.StateIdle)
		return "", err
	}

	return validName, nil
}
