package translator

import (
	"errors"
	"fmt"
	"github.com/awesome-gocui/gocui"
	"github.com/soyarielruiz/tdl-borbotones-go/client/hand"
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
		return checkDropCommand(words)
	case string(tools.EXIT):
		return checkExitCommand(words)
	case string(tools.TAKE):
		return checkTakeCommand(words)
	case "list":
		return checkListCommand(words,gui)
	default:
		return tools.Action{}, errors.New("string: Command not recognized")
	}
}

func checkDropCommand(words []string) (tools.Action, error){
	if len(words)>3{
		return hand.DropACard(words)
	}else{
		return tools.Action{},errors.New("string: Command not recognized")
	}
}

func checkTakeCommand(words []string) (tools.Action, error){
	if len(words)==1{
		return tools.Action{Command: tools.TAKE}, nil
	}else{
		return tools.Action{},errors.New("string: Command not recognized")
	}
}

func checkExitCommand(words []string) (tools.Action, error){
	if len(words)==1{
		return tools.Action{Command: tools.EXIT}, nil
	}else{
		return tools.Action{},errors.New("string: Command not recognized")
	}
}

func checkListCommand(words []string,gui *gocui.Gui) (tools.Action, error){
	if len(words)==1{
		return tools.Action{}, hand.ShowHand(gui)
	}else{
		return tools.Action{},errors.New("string: Command not recognized")
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

func showDropAction(playerId string, card tools.Card) string {
	return fmt.Sprintf("%s throws %s %s", playerId, strings.ToUpper(string(card.Suit)), strconv.Itoa(card.Number))
}

func showTakeAction(playerId string) string {
	return fmt.Sprintf("%s takes 1 card", playerId)
}

func showTurnAssigned(playerId string) string {
	return fmt.Sprintf("%s It's your turn! Drop one of your cards or take one",playerId)
}

func showExitAction(playerId string) string {
	return fmt.Sprintf("%s has left the room", playerId)
}

func MustLeave(action tools.Action) bool {
	if tools.EXIT == action.Command {
		return true
	}
	return false
}

func GameWasEnded(action tools.Action) bool {
	if tools.GAME_ENDED == action.Command {
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