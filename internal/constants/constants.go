package constants

import "github.com/go-telegram/fsm"

const (
	StateIdle       fsm.StateID = "find_clinic"
	StateStart      fsm.StateID = "start"
	StateChangeCity fsm.StateID = "change_city"
)
