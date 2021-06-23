package tools

import "fmt"

type Action struct {
	Command  Command `json:"command"`
	Card     Card    `json:"card"`
	PlayerId string  `json:"player_id"`
	Message  string  `json:"message"`
}

func (a Action) String() string {
	return fmt.Sprintf("Command:\"%s\"; Card:{\"%s\"}; PlayerId:\"%s\"; Message:\"%s\"", a.Command, a.Card, a.PlayerId, a.Message)
}

type Command string
type Suit string

type Card struct {
	Number int  `json:"number"`
	Suit   Suit `json:"suit"`
}

func (a Card) String() string {
	return fmt.Sprintf("Number:%d; Suit:\"%s\"", a.Number, a.Suit)
}

func (a Command) String() string {
	return fmt.Sprintf("Command:\"%s\"", string(a))
}

const (
	DROP Command = "drop"
	TAKE Command = "take"
	EXIT Command = "exit"
)

const (
	GREEN  Suit = "green"
	YELLOW Suit = "yellow"
	RED    Suit = "red"
	BLUE   Suit = "blue"
)

func CreateFromMessage(playerId, message string) Action {
	return Action{"", Card{}, playerId, message}
}
