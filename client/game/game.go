package game

import (
	"encoding/json"
	"log"

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

	go ga.receivingData(ga.G)
}

func (ga *Game) receivingData(g *gocui.Gui) {
	for {
		if err := ReceiveMsgFromGame(g, ga); err != nil {
			log.Fatalln(err)
		}
	}
}
