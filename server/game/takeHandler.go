package game

import (
	"github.com/soyarielruiz/tdl-borbotones-go/tools"
)

type TakeHandler struct{}

func (t TakeHandler) Handle(action tools.Action, game *Game) {
	if !game.IsUserTurn(action.PlayerId) {
		game.Users[action.PlayerId].SendChannel <- tools.CreateFromMessage(action.PlayerId, "No es tu turno!")
	} else {
		game.Users[action.PlayerId].SendChannel <- tools.Action{"", game.Deck.GetCardFromDeck(), action.PlayerId, "", []tools.Card{}}
	}
}
