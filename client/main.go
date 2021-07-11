package main

import (
	"fmt"
	"log"

	"github.com/awesome-gocui/gocui"
	"github.com/soyarielruiz/tdl-borbotones-go/client/lobby"
)

func main() {

	g, err := gocui.NewGui(gocui.OutputNormal, false)
	if err != nil {
		log.Panicln(err)
	}
	g.Cursor = true
	g.Mouse = true
	defer g.Close()

	l, err := lobby.New(g)
	if err == nil {
		g.SetManagerFunc(l.Layout)

		if err := l.Keybindings(g); err != nil {
			log.Panicln(err)
		}

		if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
			log.Panicln(err)
		}
	} else {
		g.Close()
		fmt.Println("The server is offline")
	}

}
