package hand

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
	"github.com/soyarielruiz/tdl-borbotones-go/client/translator"
	"github.com/soyarielruiz/tdl-borbotones-go/tools"
	"time"
)

var userCards []tools.Card

func ManageHand(action tools.Action) func(gui *gocui.Gui) error {
	return func(gui *gocui.Gui) error {
		err := createOrUpdateHand(gui, action)
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

func createOrUpdateHand(gui *gocui.Gui, action tools.Action) error {
	time.Sleep(1 * time.Second)
	if len(action.Cards) > 0 || len(userCards) > 0 {
		out, _ := gui.View("mano")

		if len(action.Cards) > 0 && action.Command == "" {
			userCards = action.Cards
		}

		if action.Command == tools.TAKE && action.Card.Suit != "" {
			userCards = append(userCards, action.Card)
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

func showFromServer(gui *gocui.Gui, action tools.Action) error {
	if len(action.Command.String()) > 1 {
		message, viewToUse, err := translator.TranslateMessageFromServer(action)
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