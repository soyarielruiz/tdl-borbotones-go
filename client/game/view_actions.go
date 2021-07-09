package game

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/awesome-gocui/gocui"
	"github.com/soyarielruiz/tdl-borbotones-go/client/translator"
	"github.com/soyarielruiz/tdl-borbotones-go/tools"
)

const (
	RED    = "\033[1;31m%s\033[0m"
	GREEN  = "\033[1;32m%s\033[0m"
	BLUE   = "\033[1;36m%s\033[0m"
	YELLOW = "\033[1;33m%s\033[0m"
)

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func ReceiveMsgFromGame(gui *gocui.Gui, game *Game) error {
	var err error
	for err == nil {
		var action tools.Action
		err = game.Decoder.Decode(&action)
		if err == nil {
			gui.Update(translator.ManageHand(action))
		}
		time.Sleep(1 * time.Second)
	}
	return err
}

// Layout Creates view with a division in full size
func Layout(g *gocui.Gui) error {

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen

	maxX, maxY := g.Size()
	if v, err := g.SetView("jugador", 0, 0, maxX/2-1, maxY/4-1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = "Player"
		v.Editable = true
		v.Wrap = true

		if _, err = setCurrentViewOnTop(g, "jugador"); err != nil {
			return err
		}
	}
	if v, err := g.SetView("hand", 0, maxY/4, maxX/2-1, maxY/2-1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = "Hand"
		v.Wrap = true
		v.Autoscroll = true
		v.Frame = true
	}

	if v, err := g.SetView("gamelog", 0, maxY/2, maxX/2-1, maxY-1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = "Game log"
		v.Wrap = true
		v.Autoscroll = true
	}

	if v, err := g.SetView("mesa", maxX/2-1, 0, maxX-1, maxY/2-1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = "Table"
		v.Wrap = true
		v.Autoscroll = true
	}

	if v, err := g.SetView("help", maxX/2-1, maxY/2, maxX-1, maxY-1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = "Help"
		v.Wrap = true
		v.Autoscroll = false
		v.Frame = true
		initHelp(g, v)
	}

	return nil
}

func initHelp(g *gocui.Gui, v *gocui.View) {
	fmt.Fprintf(v, RED, "======================================================================\n")
	fmt.Fprintf(v, "                           WELCOME TO ")
	fmt.Fprintf(v, RED, "G")
	fmt.Fprintf(v, BLUE, "U")
	fmt.Fprintf(v, GREEN, "N")
	fmt.Fprintf(v, YELLOW, "O\n")
	fmt.Fprintf(v, RED, "======================================================================\n\n")
	fmt.Fprintf(v, GREEN, "Game rules:\n")
	fmt.Fprintf(v, "- If it's your turn drop a card that matches the number or suit of the"+
		"  card on the table.\n")
	fmt.Fprintf(v, "- If you don't have one you must take one from the draw pile.\n")
	fmt.Fprintf(v, "- If it's not your turn but you have a card that matches the one on"+
		"     the table try to drop it and get ahead of the other players.\n\n")
	fmt.Fprintf(v, YELLOW, "Commands:\n")
	fmt.Fprintf(v, "- drop [color] [number] (e.g. drop red 5)\n"+
		"- take (takes one card from draw pile)\n- list (displays your hand)\n- exit")
}

func Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func InitKeybindings(g *gocui.Gui, ga *Game) error {

	// Bind enter key to input to send new messages.
	if err := g.SetKeybinding("jugador", gocui.KeyEnter, gocui.ModNone, jugadorBind(ga)); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, Quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyEsc, gocui.ModNone, Quit); err != nil {
		log.Panicln(err)
	}

	return nil
}

func jugadorBind(game *Game) func(g *gocui.Gui, iv *gocui.View) error {
	return func(g *gocui.Gui, iv *gocui.View) error {
		iv.Autoscroll = true
		iv.Editable = true
		// Read buffer from the beginning.
		iv.Rewind()

		// Send message if text was entered.
		if len(iv.Buffer()) >= 2 {

			messageToUse := string(iv.Buffer())
			messageToSend, err := translator.CreateAnAction(messageToUse, g)
			if err != nil {
				out, _ := g.View("gamelog")
				fmt.Fprintf(out, "Error: %s. Try again!\n", err)
			}

			if translator.MustLeave(messageToSend) {
				return Quit(g, iv)
			}
      
			if translator.HaveActionToSend(messageToSend) {
				if err := game.Encoder.Encode(&messageToSend); err != nil {
					return Quit(g, iv)
				}
			}

			//get cursor position
			_, y := iv.Cursor()

			//adding a visual enter
			w := y + 1

			err = iv.SetCursorUnrestricted(0, w)
			if err != nil {
				log.Println("Failed to set cursor:", err)
			}

			iv.Clear()
			return err
		}
		return nil
	}
}
