package lobby

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/awesome-gocui/gocui"
)

const (
	serverAddress = "127.0.0.1"
	serverPort    = "8080"
	serverConn    = "tcp"
)

type Lobby struct {
	G       *gocui.Gui
	Conn    *net.TCPConn
	Encoder *json.Encoder
	Decoder *json.Decoder
}

type LobbyOption struct {
	Option []int `json:"option"`
}

func (l *Lobby) Layout(g *gocui.Gui) error {
	l.Home(g, nil)
	return nil
}

func (l *Lobby) Keybindings(g *gocui.Gui) error {

	if err := g.SetKeybinding("home", gocui.KeyF1, gocui.ModNone, l.NewGame); err != nil {
		return err
	}

	if err := g.SetKeybinding("home", gocui.KeyF2, gocui.ModNone, l.FindGame); err != nil {
		return err
	}

	if err := g.SetKeybinding("home", gocui.KeyF3, gocui.ModNone, l.Quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("wait", gocui.KeyF3, gocui.ModNone, l.Back); err != nil {
		return err
	}

	if err := g.SetKeybinding("gamelist", gocui.KeyF3, gocui.ModNone, l.Back); err != nil {
		return err
	}

	for _, v := range "123456789" {
		if err := g.SetKeybinding("gamelist", v, gocui.ModNone, l.keyHandler(v)); err != nil {
			return err
		}
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, l.Quit); err != nil {
		return err
	}

	return nil
}

func (l *Lobby) keyHandler(ch rune) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		lo := LobbyOption{[]int{int(ch - 48)}}
		l.Encoder.Encode(lo)
		if v, err := g.SetView("selgame", 27+1, 1+1, 52+1, 3+1, gocui.TOP); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			g.SetCurrentView(v.Name())
			v.Title = "Esperando jugadores"
			fmt.Fprintln(v, ch)
		}
		go l.waitPlayer()
		return nil
	}
}

func (l *Lobby) Back(g *gocui.Gui, v *gocui.View) error {
	if g.CurrentView().Name() == "wait" {
		l.reconect()
		if err := g.DeleteView("wait"); err != nil {
			return err
		}
	}

	if g.CurrentView().Name() == "gamelist" {
		l.reconect()
		if err := g.DeleteView("gamelist"); err != nil {
			return err
		}
	}

	if _, err := g.SetCurrentView("home"); err != nil {
		return err
	}

	return nil
}

func (l *Lobby) Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (l *Lobby) Home(g *gocui.Gui, v *gocui.View) error {
	if v, err := g.SetView("home", 1, 1, 25, 5, gocui.TOP); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		g.SetCurrentView(v.Name())
		v.Title = "Bienvenidos a GUNO"
		fmt.Fprintln(v, "F1 - Nuevo Juego")
		fmt.Fprintln(v, "F2 - Buscar Partida")
		fmt.Fprintln(v, "F3 - Salir")
	}
	return nil
}

func (l *Lobby) NewGame(g *gocui.Gui, v *gocui.View) error {
	lo := LobbyOption{[]int{1}}
	l.Encoder.Encode(lo)
	if v, err := g.SetView("wait", 27, 1, 52, 3, gocui.TOP); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		g.SetCurrentView(v.Name())
		v.Title = "Esperando jugadores"
		fmt.Fprintln(v, "F3 - Atras")
	}
	go l.waitPlayer()
	return nil
}

func (l *Lobby) FindGame(g *gocui.Gui, v *gocui.View) error {
	lo := LobbyOption{[]int{2}}
	l.Encoder.Encode(lo)
	l.Decoder.Decode(&lo)
	_, y := g.Size()
	if v, err := g.SetView("gamelist", 27, 1, 52, y-3, gocui.TOP); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		g.SetCurrentView(v.Name())
		v.Title = "Lista de Partidas"
		fmt.Fprintln(v, "F3 - Atras")
		fmt.Fprintln(v, "-----------------------")
		for i, game := range lo.Option {
			fmt.Fprintln(v, strconv.Itoa(i+1)+" - Partida "+strconv.Itoa(game))
		}
	}
	return nil
}

func (l *Lobby) waitPlayer() {
	var start LobbyOption
	if err := l.Decoder.Decode(&start); err == nil {
		l.G.Update(l.Game)
	}
}

func (l *Lobby) Game(g *gocui.Gui) error {
	if v, err := g.SetView("game", 1, 1, 50, 50, gocui.TOP); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		g.SetCurrentView(v.Name())
		v.Title = "GUNO"
		fmt.Fprintln(v, "*******************************")
		fmt.Fprintln(v, "JUUUUUUUGAAAAAAAR")
	}
	return nil
}

func (l *Lobby) reconect() {
	l.Conn.Close()
	time.Sleep(1 * time.Second)
	tcpAddr, err := net.ResolveTCPAddr(serverConn, serverAddress+":"+serverPort)
	if err != nil {
		log.Panic(err)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Panic(err)
	}
	conn.SetWriteBuffer(10)
	l.Conn = conn
	l.Encoder = json.NewEncoder(conn)
	l.Decoder = json.NewDecoder(conn)
}
