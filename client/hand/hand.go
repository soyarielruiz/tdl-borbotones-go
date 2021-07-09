package hand

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/awesome-gocui/gocui"
	"github.com/soyarielruiz/tdl-borbotones-go/tools"
)

var userCards []tools.Card
var cardOnTable tools.Card

func CreateOrUpdateHand(gui *gocui.Gui, action tools.Action) error {
	if len(action.Cards) > 0 || len(userCards) > 0 {
		out, _ := gui.View("mano")

		if len(action.Cards) > 0 && action.Command == "" {
			userCards = action.Cards
			displayInitialCard(gui, action.Card)
			_ = SaveCardOnTable(action.Card)
		}

		if action.Command == tools.TAKE && action.Card.Suit != "" {
			userCards = append(userCards, action.Card)
		}

		hand := displayCards(userCards)
		_, err := fmt.Fprint(out, hand)
		if err != nil {
			return err
		}
		hand = ""
	}
	return nil
}

func displayInitialCard(gui *gocui.Gui, card tools.Card) {
	out, _ := gui.View("mesa")
	cardToShow := string(card.Suit) + " " + strconv.Itoa(card.Number)
	fmt.Fprintf(out, "Initial card: "+cardToShow+"\n")
}

func displayCards(hand []tools.Card) string {
	var handToShow []string
	for _, card := range hand {
		cardToShow := string(card.Suit) + " " + strconv.Itoa(card.Number)
		handToShow = append(handToShow, cardToShow)
	}
	return "Your cards are: " + strings.Join(handToShow, ", ") + "\n"
}

func DropACard(words []string) (tools.Action, error) {
	card := getCardFromMessage(words[1], words[2])
	_, err := itsAPlayingCard(card)

	if err != nil {
		return tools.Action{}, err
	}
	return tools.Action{Command: tools.DROP, Card: card, Cards: []tools.Card{}}, nil
}

func itsAPlayingCard(cardSent tools.Card) (interface{}, error) {
	existingPosition := -1
	for i, card := range userCards {
		if cardSent.Suit == card.Suit && cardSent.Number == card.Number {
			if cardSent.Suit == cardOnTable.Suit || cardSent.Number == cardOnTable.Number {
				existingPosition = i
			}
		}
	}

	if existingPosition == -1 {
		return false, errors.New("card: You cannot drop that card")
	}

	userCards = append(userCards[:existingPosition], userCards[existingPosition+1:]...)
	return true, nil
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
		colorToUse = ""
	}

	value, err := strconv.ParseInt(number, 10, 64)
	if err != nil || colorToUse == "" {
		return tools.Card{}
	}

	return tools.Card{Number: int(value), Suit: colorToUse}
}

func ShowHand(gui *gocui.Gui) error {
	out, _ := gui.View("mano")

	if len(userCards) > 0 {
		hand := displayCards(userCards)
		_, err := fmt.Fprint(out, hand)
		if err != nil {
			return err
		}
	}

	return nil
}

func SaveCardOnTable(card tools.Card) error {
	if cardOnTable.Suit != card.Suit || cardOnTable.Number != card.Number {
		cardOnTable = card
		return nil
	}
	return errors.New("bool: Duplicated card")
}
