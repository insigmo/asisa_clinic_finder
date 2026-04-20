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
			StateStart:      sm.CallbackStart,
			StateChangeCity: sm.CallbackChangeCity,
		},
	)

	return sm
}
