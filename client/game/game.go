package game

import (
	"encoding/json"
	"log"
	"time"
	"fmt"

	"github.com/awesome-gocui/gocui"
)

type Game struct {
	G       *gocui.Gui
	Encoder *json.Encoder
	Decoder *json.Decoder
}

func NewGame(g *gocui.Gui, enc *json.Encoder, dec *json.Decoder) *Game {
	return &Game{g, enc, dec}
}

func (ga *Game) Run() {
	
	ga.G.SetManagerFunc(Layout)
	if err := InitKeybindings(ga.G, ga); err != nil {
		log.Fatalln(err)
	}

	time.Sleep(1 * time.Second)
	ReceiveMsgFromGame(ga.G, ga)

	ga.G.Update(close())

	time.Sleep(5 * time.Second)

	ga.G.Update(func(g *gocui.Gui) error {
		return gocui.ErrQuit
	})
}

func close() func(g *gocui.Gui) error {
	return func(g *gocui.Gui) error{
		v,_ := g.View("gamelog")
		fmt.Fprintf(v,"\nConnection lost. The window will close in a few seconds...\n")
		return nil
	}
}
