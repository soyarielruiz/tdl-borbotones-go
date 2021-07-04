package hand

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
	"github.com/soyarielruiz/tdl-borbotones-go/client/translator"
	"github.com/soyarielruiz/tdl-borbotones-go/tools"
)

var userCards []tools.Card

func ManageHand(action tools.Action) func(gui *gocui.Gui) error {
	return func(gui *gocui.Gui) error {
		createOrUpdateHand(gui, action.Cards)

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
}

func createOrUpdateHand(gui *gocui.Gui, cards []tools.Card) error {
	if len(cards) > 0 || len(userCards) > 0 {
		out, _ := gui.View("mano")
		if len(cards) > 0 && len(userCards) == 0 {
			userCards = cards
		}

		if len(userCards) > 0 && len(cards) == 1 {
			userCards = append(userCards, cards...)
		}

		hand := translator.DisplayCards(userCards)
		_, err := fmt.Fprintf(out, hand)
		if err != nil {
			return err
		}
		hand = ""
	}
	return nil
}