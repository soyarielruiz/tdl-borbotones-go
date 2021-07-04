package game

import (
	"github.com/soyarielruiz/tdl-borbotones-go/tools"
)

type TakeHandler struct{}

func (t TakeHandler) Handle(action tools.Action, game *Game) {
	if !game.Tur.IsUserTurn(action.PlayerId) {
		game.Users[action.PlayerId].SendChannel <- tools.CreateFromMessage(action.PlayerId, "It's not your turn!")
	} else {
		game.Users[action.PlayerId].CardsLeft++
		game.Users[action.PlayerId].SendChannel <- tools.Action{Command: tools.TAKE, Card: game.Deck.GetCardFromDeck(), PlayerId: action.PlayerId, Cards: []tools.Card{}}
	}
}
