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
	} else if action.Card.Number == game.Deck.GetFrontCard().Number {
		game.Deck.PutCard(action.Card)
		game.SendToAll(&action)
		game.Tur.GoTo(action.PlayerId)
		game.Tur.Next()
	} else {
		game.Users[action.PlayerId].SendChannel <- tools.CreateFromMessage(action.PlayerId, "No es tu turno!")
	}
}
