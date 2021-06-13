package action

type Action struct {
	Command  Command `json:"command"`
	Card     Card    `json:"card"`
	PlayerId string  `json:"player_id"`
	Message  string  `json:"message"`
}

type Command string
type Suit string

type Card struct {
	Number int  `json:"number"`
	Suit   Suit `json:"suit"`
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
