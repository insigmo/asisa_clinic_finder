package fsm_states

import (
	"fmt"
	"log"

	"github.com/go-telegram/fsm"
	"github.com/insigmo/asisa_clinic_finder/internal/db"
	"github.com/insigmo/asisa_clinic_finder/internal/helpers"
	"github.com/insigmo/asisa_clinic_finder/internal/local_models"
	"github.com/insigmo/asisa_clinic_finder/internal/localize_manager"
)

func (s *StateMachine) CallbackStart(_ *fsm.FSM, args ...any) {
	params := args[0].(*local_models.BaseParams)
	localizer := localize_manager.New(params.Update.Message.From.LanguageCode)

	err := helpers.SendMessage(params, localizer.StartMessage())
	if err != nil {
		log.Fatal(fmt.Sprintf("error in start callback %v", err))
	}

	s.FSM.Transition(params.UserID, StateChangeCity)
}

func (s *StateMachine) CallbackChangeCity(_ *fsm.FSM, args ...any) {
	params := args[0].(*local_models.BaseParams)
	localizer := localize_manager.New(params.Update.Message.From.LanguageCode)

	err := helpers.SendMessage(params, localizer.StartMessage())
	if err != nil {
		log.Fatal(fmt.Sprintf("error in start callback %v", err))
	}

	dbManager := params.Ctx.Value("dbManager").(*db.Manager)
	userInfo := params.Update.Message.From

	// TODO add validate city

	city := params.Update.Message.Text
	user := &db.User{
		ID:           params.UserID,
		Name:         userInfo.FirstName,
		Lastname:     userInfo.LastName,
		IsBot:        userInfo.IsBot,
		City:         city,
		LanguageCode: userInfo.LanguageCode,
	}

	err = dbManager.InsertOrUpdateUser(params.Ctx, user)
	if err != nil {
		panic(err)
	}
}
