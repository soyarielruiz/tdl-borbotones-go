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
		//var initCard tools.Card

		if len(action.Cards) > 0 && action.Command == "" {
			userCards = action.Cards
			//initCard=action.Card
			hand := displayCards(userCards)
			fmt.Fprintf(out, hand)
		} else {
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
	}
	return nil
}

func displayCards(hand []tools.Card) string {
	var handToShow []string
	for _, card := range hand {
		cardToShow := string(card.Suit) + " " + strconv.Itoa(card.Number)
		handToShow = append(handToShow, cardToShow)
	}
	return "Your cards are: " + strings.Join(handToShow, ", ") + "\n"
}

func displayInitialCard(gui *gocui.Gui, initCard tools.Card){
	initialCard:=string(initCard.Suit) + " " + strconv.Itoa(initCard.Number)
	out, _ := gui.View("mesa")
	fmt.Fprintf(out,"Initial card: " +initialCard + "\n")
}

func DropACard(words []string) (tools.Action, error) {
	card := getCardFromMessage(words[1], words[2])
	_, err := itsAPlayingCard(card)
	if err != nil {
		return tools.Action{}, err
	}
	message := strings.Join(words[3:], " ")
	return tools.Action{Command: tools.DROP, Card: card, Message: message, Cards: []tools.Card{}}, nil
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