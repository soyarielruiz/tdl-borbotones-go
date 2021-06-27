package game

import (
	"github.com/soyarielruiz/tdl-borbotones-go/tools"
)

type DropHandler struct{}

func (t DropHandler) Handle(action tools.Action, game *Game) {
	if game.IsUserTurn(action.PlayerId) {
		game.Deck.PutCard(action.Card)
		game.SendToAll(&action)
		game.Turnero.Next()
	} else if action.Card.Number == game.Deck.GetFrontCard().Number {
		game.Deck.PutCard(action.Card)
		game.SendToAll(&tools.Action{})
		game.Turnero.GoTo(action.PlayerId)
		game.Turnero.Next()
	} else {
		game.Users[action.PlayerId].SendChannel <- tools.CreateFromMessage(action.PlayerId, "No es tu turno!")
	}
}
