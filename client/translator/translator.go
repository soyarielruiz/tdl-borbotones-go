package translator

import (
	"fmt"
	"github.com/soyarielruiz/tdl-borbotones-go/tools"
	"os"
	"strconv"
	"strings"
)

func CreateAnAction(messageToSend string) tools.Action {
	words := strings.Fields(messageToSend)

	command := getCommandFromMessage(words[0])
	card := getCardFromMessage(words[1], words[2])
	message := strings.Join(words[3:], " ")

	return tools.Action{command, card, "", message}
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

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}