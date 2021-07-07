package game

import (
	"encoding/json"
	"log"
	"time"

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
	go ReceiveMsgFromGame(ga.G, ga)
}
