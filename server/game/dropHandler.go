package game

import (
	"github.com/soyarielruiz/tdl-borbotones-go/tools"
)

type DropHandler struct{}

func (t DropHandler) Handle(action tools.Action, game *Game) {
	if !game.IsUserTurn(action.PlayerId) {
		game.Users[action.PlayerId].SendChannel <- tools.CreateFromMessage(action.PlayerId, "No es tu turno!")
	} else {
		//TODO
	}
}
