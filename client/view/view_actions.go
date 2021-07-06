package view

import (
	"errors"
	"github.com/awesome-gocui/gocui"
	"github.com/soyarielruiz/tdl-borbotones-go/client/translator"
	"github.com/soyarielruiz/tdl-borbotones-go/client/game"
	"github.com/soyarielruiz/tdl-borbotones-go/tools"
	"log"
	"time"
	"fmt"
	"os"
)

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func ReceiveMsgFromGame(gui *gocui.Gui, game *game.Game) error {
	//wait a starting moment
	time.Sleep(1*time.Second)
	for {
		var action tools.Action
		err:=game.Decoder.Decode(&action)
		if err !=nil{
			fmt.Fprintf(os.Stderr, "Error en decode : %s\n", err)
			if "EOF" == err.Error() {
				fmt.Fprintf(os.Stderr, "Fatal error in conection: %s ", err.Error())
				os.Exit(1)
			}
		}
		go gui.Update(translator.ManageHand(action))
	}
}

// Layout Creates view with a division in full size
func Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("jugador", 0, 0, maxX/2-1, maxY/2-1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = "Jugador"
		v.Editable = true
		v.Wrap = true

		if _, err = setCurrentViewOnTop(g, "jugador"); err != nil {
			return err
		}
	}

	if v, err := g.SetView("mano", 0, maxY/2, maxX-1, maxY-1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = "Mano"
		v.Wrap = true
		v.Autoscroll = true
	}

	if v, err := g.SetView("mesa", maxX/2-1, 0, maxX-1, maxY/2-1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = "Mesa"
		v.Wrap = true
		v.Autoscroll = true
	}

	return nil
}

func Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func InitKeybindings(g *gocui.Gui) error {

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, Quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyEsc, gocui.ModNone, Quit); err != nil {
		log.Panicln(err)
	}

	return nil
}