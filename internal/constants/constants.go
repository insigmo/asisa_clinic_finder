package constants

type StateID string

const (
	StateIdle       StateID = "find_clinic"
	StateStart      StateID = "start"
	StateChangeCity StateID = "change_city"
)
