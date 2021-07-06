package game

import (
	"fmt"
	"github.com/soyarielruiz/tdl-borbotones-go/tools"
	"strconv"
)

type DropHandler struct{}

func (t DropHandler) Handle(action tools.Action, game *Game) {
	frontCard := game.Deck.GetFrontCard()
	if game.Tur.IsUserTurn(action.PlayerId) {
		game.Users[action.PlayerId].SendChannel <- tools.CreateFromMessage(action.PlayerId, fmt.Sprintf("CARTA EN MESA %s %s", strconv.Itoa(frontCard.Number), frontCard.Suit))
		if action.Card.Number == frontCard.Number || action.Card.Suit == frontCard.Suit {
			game.Deck.PutCard(action.Card)
			game.SendToAll(&action)
			game.TurnMoveForward()
			game.Users[action.PlayerId].CardsLeft--
			if game.Users[action.PlayerId].CardsLeft == 0 {
				game.Ended = true
				game.SendToAll(&tools.Action{
					Command:  tools.GAME_ENDED,
					Card:     tools.Card{},
					PlayerId: action.PlayerId,
					Message:  "Game ended! Player " + action.PlayerId + " won!",
					Cards:    nil,
				})
			}
		} else {
			game.Users[action.PlayerId].SendChannel <- tools.CreateFromMessage(action.PlayerId, "Jugada invalida")
		}
	} else if action.Card.Number == frontCard.Number &&
		action.Card.Suit == frontCard.Suit {
		game.Deck.PutCard(action.Card)
		game.SendToAll(&action)
		game.Tur.GoTo(action.PlayerId)
		game.TurnMoveForward()
	} else {
		game.Users[action.PlayerId].SendChannel <- tools.CreateFromMessage(action.PlayerId, "No es tu turno!")
	}
}
