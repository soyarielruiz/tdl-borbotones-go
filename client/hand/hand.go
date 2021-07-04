package hand

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
	"github.com/soyarielruiz/tdl-borbotones-go/client/translator"
	"github.com/soyarielruiz/tdl-borbotones-go/tools"
)

var userCards tools.Action

func ManageHand(action tools.Action, gui *gocui.Gui) error {

	if len(action.Cards) > 0 || len(userCards.Cards) > 0 {
		out, _ := gui.View("mano")
		if len(action.Cards) > 0 && len(userCards.Cards) < len(action.Cards) {
			userCards = action
		}
		hand := translator.DisplayCards(action)
		_, err := fmt.Fprintf(out, hand)
		if err != nil {
			return err
		}
		hand = ""
	}

	if len(action.Command.String()) > 1 {
		out, _ := gui.View("mesa")
		message, err := translator.TranslateMessageFromServer(action)
		if err == nil {
			_, err := fmt.Fprintln(out, message)
			if err != nil {
				return err
			}
			message = ""
		}
	}

	return nil
}
