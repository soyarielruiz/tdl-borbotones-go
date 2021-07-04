package view

import (
	"encoding/json"
	"errors"
	"github.com/awesome-gocui/gocui"
	"github.com/soyarielruiz/tdl-borbotones-go/client/hand"
	"github.com/soyarielruiz/tdl-borbotones-go/tools"
	"log"
	"net"
	"time"
)

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func ReceiveMsgFromGame(gui *gocui.Gui, conn *net.TCPConn) error {
	//wait a starting moment
	time.Sleep(1*time.Second)
	for {
		decoder := json.NewDecoder(conn)
		var action tools.Action
		decoder.Decode(&action)
		go gui.Update(hand.ManageHand(action))
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