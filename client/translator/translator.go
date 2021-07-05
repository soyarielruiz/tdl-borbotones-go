package translator

import (
	"errors"
	"fmt"
	"github.com/awesome-gocui/gocui"
	"github.com/soyarielruiz/tdl-borbotones-go/client/hand"
	"os"
	"strconv"
	"strings"

	"github.com/soyarielruiz/tdl-borbotones-go/tools"
)

func CreateAnAction(messageToSend string, gui *gocui.Gui) (tools.Action, error) {
	words := strings.Fields(messageToSend)
	action, err := validateCommand(words, gui)
	return action, err
}

func validateCommand(words []string, gui *gocui.Gui) (tools.Action, error) {
	if len(words) < 1 {
		return tools.Action{}, errors.New("string: Need more parameters")
	}

	action, err := createActionFromCommand(words, gui)

	if err != nil {
		return tools.Action{}, err
	}

	return action, nil
}

func createActionFromCommand(words []string, gui *gocui.Gui) (tools.Action, error) {

	switch strings.ToLower(words[0]) {
	case string(tools.DROP):
		return hand.DropACard(words)
	case string(tools.EXIT):
		return tools.Action{Command: tools.EXIT}, nil
	case string(tools.TAKE):
		return tools.Action{Command: tools.TAKE}, nil
	case "list":
		return tools.Action{}, hand.ShowHand(gui)
	default:
		return tools.Action{}, errors.New("string: Command not recognized")
	}
}

func TranslateMessageFromServer(action tools.Action) (string, string, error) {
	var response string
	var out string

	if string(action.Command) == string(tools.TURN_ASSIGNED) {
		response = showTurnAssigned(action.PlayerId[:5])
		out = "mano"
		return response, out, nil
	}

	if len(action.Command.String()) > 1 {
		switch strings.ToLower(string(action.Command)) {
		case string(tools.DROP):
			response = showDropAction(action.PlayerId[:5], action.Card)
			out = "mesa"
		case string(tools.EXIT):
			response = showExitAction(action.PlayerId[:5])
			out = "mano"
		case string(tools.TAKE):
			response = showTakeAction(action.PlayerId[:5])
			out = "mano"
		case string(tools.GAME_ENDED):
			response = "Game Finalizado"
			out = "mano"
		default:
			response = ""
			out = ""
		}

		return response, out, nil
	}
	return "", "", errors.New("object:Wrong action")
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

func showTurnAssigned(playerId string) string {
	return fmt.Sprintf("%s It's your turn! Drop one of your cards or take one",playerId)
}

func showExitAction(playerId string) string {
	return fmt.Sprintf("%s ha salido de la partida", playerId)
}

func MustLeave(action tools.Action) bool {
	if tools.EXIT == action.Command {
		return true
	}
	return false
}

func ManageHand(action tools.Action) func(gui *gocui.Gui) error {
	return func(gui *gocui.Gui) error {
		err := hand.CreateOrUpdateHand(gui, action)
		if err != nil {
			return err
		}
		err = showFromServer(gui, action)
		if err != nil {
			return err
		}
		return nil
	}
}

func showFromServer(gui *gocui.Gui, action tools.Action) error {
	if len(action.Command) > 1 {
		message, viewToUse, err := TranslateMessageFromServer(action)
		if err == nil {
			out, _ := gui.View(viewToUse)
			_, err := fmt.Fprintln(out, message)
			if err != nil {
				return err
			}
			message = ""
		}
	}
	return nil
}

func HaveActionToSend(action tools.Action) bool {
	if action.Command != "" || len(action.Command) > 0 {
		return true
	}
	return false
}