package fsmstate

import (
	"strings"

	"github.com/go-telegram/fsm"

	"github.com/insigmo/asisa_clinic_finder/internal/constants"
	"github.com/insigmo/asisa_clinic_finder/internal/db"
	"github.com/insigmo/asisa_clinic_finder/internal/helpers"
	"github.com/insigmo/asisa_clinic_finder/internal/i18n"
	"github.com/insigmo/asisa_clinic_finder/internal/model"
)

const (
	MinLanguageCodeLen = 2
)

func (s *StateMachine) CallbackStart(_ *fsm.FSM, args ...any) {
	if len(args) == 0 {
		return
	}

	params, ok := args[0].(*model.BaseParams)
	if !ok {
		return
	}

	localizer := i18n.New(params.Update.Message.From.LanguageCode)
	if err := helpers.SendMessage(params, localizer.StartMessage()); err != nil {
		params.Log.Error(err.Error())
		return
	}

	s.FSM.Transition(params.UserID, constants.StateChangeCity, params)
}

func (s *StateMachine) CallbackChangeCity(_ *fsm.FSM, args ...any) {
	if len(args) == 0 {
		return
	}

	params, _ := args[0].(*model.BaseParams)
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
		State:        "",
	}

	if err := dbManager.InsertOrUpdateUser(params.Ctx, user); err != nil {
		params.Log.Error(err.Error())

		return
	}

	localizer := i18n.New(userInfo.LanguageCode)
	if err := helpers.SendMessage(params, localizer.SaveUserMessage()); err != nil {
		params.Log.Error(err.Error())
	}
}
