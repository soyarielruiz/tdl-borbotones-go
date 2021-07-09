package game

import (
	"github.com/soyarielruiz/tdl-borbotones-go/tools"
)

type DropHandler struct{}

func (t DropHandler) Handle(action tools.Action, game *Game) {
	frontCard := game.Deck.GetFrontCard()
	if game.Tur.IsUserTurn(action.PlayerId) {
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
			game.Users[action.PlayerId].SendChannel <- tools.CreateFromMessage(action.PlayerId, "Invalid move. Card must match number or suit of the card on the table.Try again!")
		}
	} else if action.Card.Number == frontCard.Number &&
		action.Card.Suit == frontCard.Suit {
		game.Deck.PutCard(action.Card)
		game.SendToAll(&action)
		game.Tur.GoTo(action.PlayerId)
		game.TurnMoveForward()
	} else {
		game.Users[action.PlayerId].SendChannel <- tools.CreateFromMessage(action.PlayerId, "It's not your turn!")
	}
}
