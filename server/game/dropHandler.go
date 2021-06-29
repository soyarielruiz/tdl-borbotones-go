package game

import (
	"github.com/soyarielruiz/tdl-borbotones-go/tools"
)

type DropHandler struct{}

func (t DropHandler) Handle(action tools.Action, game *Game) {
	if game.Tur.IsUserTurn(action.PlayerId) {
		game.Deck.PutCard(action.Card)
		game.SendToAll(&action)
		game.Tur.Next()
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

	} else if action.Card.Number == game.Deck.GetFrontCard().Number {
		game.Deck.PutCard(action.Card)
		game.SendToAll(&action)
		game.Tur.GoTo(action.PlayerId)
		game.Tur.Next()
	} else {
		game.Users[action.PlayerId].SendChannel <- tools.CreateFromMessage(action.PlayerId, "No es tu turno!")
	}
}
