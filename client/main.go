package main

import (
	"log"

	"github.com/awesome-gocui/gocui"
	"github.com/soyarielruiz/tdl-borbotones-go/client/lobby"
)

func main() {

	g, err := gocui.NewGui(gocui.OutputNormal, false)
	g.Cursor = true
	g.Mouse = true
	l := lobby.New(g)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(l.Layout)

	if err := l.Keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}
