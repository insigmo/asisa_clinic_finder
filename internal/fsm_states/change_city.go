package fsm_states

import (
	"strings"

	"github.com/go-telegram/fsm"

	"github.com/insigmo/asisa_clinic_finder/internal/db"
	"github.com/insigmo/asisa_clinic_finder/internal/helpers"
	"github.com/insigmo/asisa_clinic_finder/internal/local_models"
	"github.com/insigmo/asisa_clinic_finder/internal/localize_manager"
)

const MinLanguageCodeLen = 2

func (s *StateMachine) CallbackStart(_ *fsm.FSM, args ...any) {
	if len(args) == 0 {
		return
	}

	params, ok := args[0].(*local_models.BaseParams)
	if !ok {
		return
	}

	localizer := localize_manager.New(params.Update.Message.From.LanguageCode)
	if err := helpers.SendMessage(params, localizer.StartMessage()); err != nil {
		params.Log.Error(err.Error())
		return
	}

	s.FSM.Transition(params.UserID, StateChangeCity, params)
}

func (s *StateMachine) CallbackChangeCity(_ *fsm.FSM, args ...any) {
	if len(args) == 0 {
		return
	}

	params, _ := args[0].(*local_models.BaseParams)
	dbManager, ok := params.Ctx.Value("dbManager").(*db.Manager)
	if !ok {
		params.Log.Error("db manager not found in context")

		return
	}

	userInfo := params.Update.Message.From
	city := strings.TrimSpace(params.Update.Message.Text)

	if len([]rune(city)) < MinLanguageCodeLen {
		return
	}

	user := &db.User{
		ID:           userInfo.ID,
		Username:     userInfo.Username,
		Name:         userInfo.FirstName,
		Lastname:     userInfo.LastName,
		IsBot:        userInfo.IsBot,
		City:         city,
		LanguageCode: userInfo.LanguageCode,
	}

	if err := dbManager.InsertOrUpdateUser(params.Ctx, user); err != nil {
		params.Log.Error(err.Error())

		return
	}

	localizer := localize_manager.New(userInfo.LanguageCode)
	if err := helpers.SendMessage(params, localizer.SaveUserMessage()); err != nil {
		params.Log.Error(err.Error())
	}
}
