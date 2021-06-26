package commandHandler

import (
	"github.com/soyarielruiz/tdl-borbotones-go/server/game"
	"github.com/soyarielruiz/tdl-borbotones-go/tools"
)

type CommandHandler interface {
	Handle(action tools.Action, game *game.Game)
}
