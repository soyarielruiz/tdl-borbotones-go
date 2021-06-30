package translator

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/soyarielruiz/tdl-borbotones-go/tools"
)

func CreateAnAction(messageToSend string) tools.Action {
	words := strings.Fields(messageToSend)

	command := getCommandFromMessage(words[0])
	card := getCardFromMessage(words[1], words[2])
	message := strings.Join(words[3:], " ")

	return tools.Action{command, card, "", message, []tools.Card{}}
}

func getCardFromMessage(color string, number string) tools.Card {
	var colorToUse = tools.GREEN
	switch strings.ToLower(color) {
	case string(tools.GREEN):
		colorToUse = tools.GREEN
	case string(tools.YELLOW):
		colorToUse = tools.YELLOW
	case string(tools.RED):
		colorToUse = tools.RED
	case string(tools.BLUE):
		colorToUse = tools.BLUE
	default:
		colorToUse = tools.GREEN
	}

	value, err := strconv.ParseInt(number, 10, 64)
	checkError(err)

	return tools.Card{int(value), colorToUse}
}

func getCommandFromMessage(message string) tools.Command {

	switch strings.ToLower(message) {
	case string(tools.DROP):
		return tools.DROP
	case string(tools.EXIT):
		return tools.EXIT
	case string(tools.TAKE):
		return tools.TAKE
	default:
		return tools.DROP
	}
}

func TranslateMessageFromServer(action tools.Action) (string, error) {
	var response string
	//strings.ToUpper(string(action.Command)) +
	if len(action.Command.String()) > 1 {
		switch strings.ToLower(string(action.Command)) {
		case string(tools.DROP):
			response = showDropAction(string(action.PlayerId)[:5], action.Card)
		case string(tools.EXIT):
			response = showExitAction(string(action.PlayerId)[:5])
		case string(tools.TAKE):
			response = showTakeAction(string(action.PlayerId)[:5])
		default:
			response = ""
		}

		return response, nil
	}
	return "", errors.New("object:Wrong action")
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func showDropAction(playerId string, card tools.Card) string {
	return fmt.Sprintf("%s lanza %s %s", playerId, strings.ToUpper(string(card.Suit)), strconv.Itoa(card.Number))
}

func showTakeAction(playerId string) string {
	return fmt.Sprintf("%s toma 1 carta", playerId)
}

func showExitAction(playerId string) string {
	return fmt.Sprintf("%s ha salido de la partida", playerId)
}

func DisplayCards(hand tools.Action) string {
	var handToShow []string
	for _, card := range hand.Cards {
		cardToShow := string(card.Suit) + " " + strconv.Itoa(card.Number)
		handToShow = append(handToShow, cardToShow)
	}
	return "Tus cartas son: " + strings.Join(handToShow, ", ")
}