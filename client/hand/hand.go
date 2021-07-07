package hand

import (
	"errors"
	"fmt"
	"github.com/awesome-gocui/gocui"
	"github.com/soyarielruiz/tdl-borbotones-go/tools"
	"os"
	"strconv"
	"strings"
	"time"
)

var userCards []tools.Card

func CreateOrUpdateHand(gui *gocui.Gui, action tools.Action) error {
	time.Sleep(1 * time.Second)
	if len(action.Cards) > 0 || len(userCards) > 0 {
		out, _ := gui.View("mano")

		if len(action.Cards) > 0 && action.Command == "" {
			userCards = action.Cards
			displayInitialCard(gui,action.Card)
		}

		if action.Command == tools.TAKE && action.Card.Suit != "" {
			userCards = append(userCards, action.Card)
		}

		hand := displayCards(userCards)
		_, err := fmt.Fprintf(out, hand)
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
	fmt.Fprintf(out, "Initial card: " + cardToShow + "\n")
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
			existingPosition = i
		}
	}

	if existingPosition == -1 {
		return false, errors.New("card: No posees esa carta en tu mano")
	}

	userCards = append(userCards[:existingPosition], userCards[existingPosition+1:]...)
	return true, nil
}

func getCardFromMessage(color string, number string) tools.Card{
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
	if err!=nil || colorToUse == ""{
		return tools.Card{}
	}

	return tools.Card{Number: int(value), Suit: colorToUse}
}

func ShowHand(gui *gocui.Gui) (error) {
	out, _ := gui.View("mano")

	if len(userCards) > 0 {
		hand := displayCards(userCards)
		_, err := fmt.Fprintf(out, hand)
		if err != nil {
			return err
		}
	}

	return nil
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}