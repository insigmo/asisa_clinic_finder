package fsm_states

import (
	"github.com/go-telegram/fsm"
)

type StateMachine struct {
	FSM *fsm.FSM
}

const (
	StateStart      fsm.StateID = "start"
	StateChangeCity fsm.StateID = "change_city"
)

func New() *StateMachine {
	sm := new(StateMachine)
	sm.FSM = fsm.New(
		StateStart,
		map[fsm.StateID]fsm.Callback{
			"stateStart":      sm.CallbackStart,
			"stateChangeCity": sm.CallbackChangeCity,
		},
	)

	return sm
}
