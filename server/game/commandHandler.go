package game

import (
	"github.com/soyarielruiz/tdl-borbotones-go/tools"
)

type CommandHandler interface {
	Handle(action tools.Action, game *Game)
}
