package handlers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/insigmo/asisa_clinic_finder/internal/helpers"
	"github.com/insigmo/asisa_clinic_finder/internal/local_models"
)

func handleDefault(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	params := local_models.NewBaseParams(ctx, tgBot, update)

	switch sm.FSM.Current(params.UserID) {
	case StateStart:
		helpers.SendMessage(params, "Send /form to start the form.")

	case StateChangeCity:
		name := update.Message.Text
		if len([]rune(name)) < 2 {
			helpers.SendMessage(params, "Name is too short, please enter at least 2 characters.")
			return
		}

		app.f.Set(userID, keyName, name)

		// onAskAge callback will prompt the user.
		app.f.Transition(userID, stateAskAge, ctx, chatID)

	case stateAskAge:
		age, err := strconv.Atoi(update.Message.Text)
		if err != nil || age < 18 || age > 100 {
			send(ctx, b, chatID, "Please enter a valid age between 18 and 100.")
			return
		}

		app.f.Set(userID, keyAge, update.Message.Text)

		// onFinish callback will display the summary and reset the state.
		app.f.Transition(userID, stateFinish, ctx, chatID, userID)

	default:
		fmt.Printf("unexpected state: %s\n", app.f.Current(userID))
	}
}
