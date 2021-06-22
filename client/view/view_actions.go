package view

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/awesome-gocui/gocui"
	"io"
	"log"
	"net"
)


func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

//TODO: Receive message from game and print it
func ReceiveMsgFromGame(g *gocui.Gui, v *gocui.View, conn *net.TCPConn) error {
	// receives msg from server
	var messageFromServer bytes.Buffer
	io.Copy(&messageFromServer, conn)
	name := "Mesa"
	g.Cursor = false

	out, err := g.View("mesa")
	if err != nil {
		return err
	}
	fmt.Fprintln(out, "Server: %s", messageFromServer.String())

	if _, err := setCurrentViewOnTop(g, name); err != nil {
		return err
	}

	return nil
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

	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Panicln(err)
	}

	return nil
}