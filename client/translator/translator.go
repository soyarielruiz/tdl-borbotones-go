package translator

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"github.com/soyarielruiz/tdl-borbotones-go/tools/action"
)

func CreateAnAction(messageToSend string) action.Action {
	words := strings.Fields(messageToSend)

	command := getCommandFromMessage(words[0])
	card := getCardFromMessage(words[1], words[2])
	message := strings.Join(words[3:], " ")

	return action.Action{command, card, "", message}
}

func getCardFromMessage(color string, number string) action.Card {
	var colorToUse = action.GREEN
	switch strings.ToLower(color) {
	case action.GREEN:
		colorToUse = action.GREEN
		break
	case action.YELLOW:
		colorToUse = action.YELLOW
		break
	case action.RED:
		colorToUse = action.RED
		break
	case action.BLUE:
		colorToUse = action.BLUE
		break
	default:
		colorToUse = action.GREEN
	}

	value, err := strconv.ParseInt(number, 10, 64)
	checkError(err)

	return action.Card{value, colorToUse}
}

func getCommandFromMessage(message string) action.Command {

	switch strings.ToLower(message) {
	case action.DROP:
		return action.DROP
	case action.EXIT:
		return action.EXIT
	case action.TAKE:
		return action.TAKE
	default:
		return action.DROP
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func translateMessage(messageToTranslate []byte) (error, Message) {
	var msg Message

	err := json.Unmarshal(messageToTranslate, &msg)

	if err != nil {
		log.Fatalf("Message cannot be decoding")
	}

	return err, msg
}

func SendMessage(messageToSend string) string {
	messageJson := ToJSON(messageToSend)
	return messageJson
}

func ReceiveMessage(messageToReceive []byte) interface{} {
	_, messageToString := translateMessage(messageToReceive)
	return messageToString
}