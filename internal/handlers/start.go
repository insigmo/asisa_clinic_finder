package handlers

import (
	"context"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/insigmo/asisa_clinic_finder/internal/fsm_states"
	"github.com/insigmo/asisa_clinic_finder/internal/helpers"
	"github.com/insigmo/asisa_clinic_finder/internal/local_models"
	"github.com/insigmo/asisa_clinic_finder/internal/localize_manager"
)

func Start(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	params := local_models.NewBaseParams(ctx, tgBot, update)
	localizator := localize_manager.New(update.Message.From.LanguageCode)

	msg := localizator.StartMessage()
	stateMachine := fsm_states.New()

	err := helpers.SendMessage(params, msg)

	if err != nil {
		log.Fatal(err)
	}

	params.Update = update
	stateMachine.FSM.Transition(params.UserID, fsm_states.StateChangeCity, params)
}
