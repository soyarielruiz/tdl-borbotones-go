package tools

import "fmt"

type Action struct {
	Command  Command `json:"command"`
	Card     Card    `json:"card"`
	PlayerId string  `json:"player_id"`
	Message  string  `json:"message"`
	Cards    []Card  `json:"cards"`
}

func (a Action) String() string {
	return fmt.Sprintf("Command:\"%s\"; Card:{Number:%d, Suit: %s}; PlayerId:\"%s\"; Message:\"%s\";Cards :%s", a.Command,
		a.Card.Number, a.Card.Suit, a.PlayerId, a.Message, a.Cards)
}

type Command string
type Suit string

type Card struct {
	Number int  `json:"number"`
	Suit   Suit `json:"suit"`
}

func (a Card) String() (string, string) {
	return string(rune(a.Number)), string(a.Suit)
}

func (a Command) String() string {
	return string(a)
}

const (
	DROP          Command = "drop"
	TAKE          Command = "take"
	EXIT          Command = "exit"
	GAME_ENDED    Command = "game_ended"
	TURN_ASSIGNED Command = "turn_assigned"
)

const (
	GREEN  Suit = "green"
	YELLOW Suit = "yellow"
	RED    Suit = "red"
	BLUE   Suit = "blue"
)

func Suits() []Suit {
	return []Suit{GREEN, YELLOW, RED, BLUE}
}

func CreateFromMessage(playerId, message string) Action {
	return Action{"", Card{}, playerId, message, []Card{}}
}
