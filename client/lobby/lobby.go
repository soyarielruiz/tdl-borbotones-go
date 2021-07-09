package lobby

import (
	"encoding/json"
	"fmt"
	"github.com/awesome-gocui/gocui"
	"github.com/soyarielruiz/tdl-borbotones-go/client/game"
	"log"
	"net"
	"strconv"
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
	Games   []int
}

type LobbyOption struct {
	Option   []int  `json:"option"`
	Nickname string `json:"nickname"`
}

func New(g *gocui.Gui) *Lobby {
	tcpAddr, _ := net.ResolveTCPAddr(serverConn, serverAddress+":"+serverPort)
	conn, _ := net.DialTCP("tcp", nil, tcpAddr)
	conn.SetWriteBuffer(10)
	return &Lobby{g, conn, json.NewEncoder(conn), json.NewDecoder(conn), nil}
}

func (l *Lobby) Layout(g *gocui.Gui) error {
	l.Home(g, nil)
	return nil
}

func (l *Lobby) Keybindings(g *gocui.Gui) error {

	if err := g.SetKeybinding("nick", gocui.KeyF1, gocui.ModNone, l.NewGame); err != nil {
		return err
	}

	if err := g.SetKeybinding("nick", gocui.KeyF2, gocui.ModNone, l.FindGame); err != nil {
		return err
	}

	if err := g.SetKeybinding("nick", gocui.KeyF3, gocui.ModNone, l.Quit); err != nil {
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
		key := int(ch - 48)
		if 0 < key && key <= len(l.Games) {
			lo := LobbyOption{[]int{key}, ""}
			l.Encoder.Encode(lo)
			if err := g.DeleteView("gamelist"); err != nil {
				return err
			}
			if v, err := g.SetView("wait", 27, 1, 52, 4, gocui.TOP); err != nil {
				if err != gocui.ErrUnknownView {
					return err
				}
				g.SetCurrentView(v.Name())
				v.Title = "Esperando jugadores"
				fmt.Fprintln(v, "F3 - Atras")
				fmt.Fprintf(v, "Opcion seleccionada %d\n", key)
			}
			go l.waitPlayer()
		}
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

	if _, err := g.SetCurrentView("nick"); err != nil {
		return err
	}

	return nil
}

func (l *Lobby) Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (l *Lobby) Home(g *gocui.Gui, v *gocui.View) error {
	if v, err := g.SetView("menu", 1, 1, 25, 5, gocui.TOP); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Bienvenidos a GUNO"
		fmt.Fprintln(v, "F1 - Nuevo Juego")
		fmt.Fprintln(v, "F2 - Buscar Partida")
		fmt.Fprintln(v, "F3 - Salir")
	}

	if v, err := g.SetView("nick", 1, 6, 25, 8, gocui.TOP); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		g.SetCurrentView(v.Name())
		v.Title = "Nick"
		v.Editable = true
		v.Wrap = true
	}

	return nil
}

func (l *Lobby) NewGame(g *gocui.Gui, v *gocui.View) error {
	nick := v.Buffer()
	if nick != "" {
		lo := LobbyOption{[]int{1}, v.Buffer()}
		l.Encoder.Encode(lo)
		if v, err := g.SetView("wait", 27, 1, 52, 4, gocui.TOP); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			g.SetCurrentView(v.Name())
			v.Title = "Esperando jugadores"
			fmt.Fprintln(v, "F3 - Atras")
			fmt.Fprintln(v, "Nueva Partida")
		}
		go l.waitPlayer()
	}
	return nil
}

func (l *Lobby) FindGame(g *gocui.Gui, v *gocui.View) error {
	nick := v.Buffer()
	if nick != "" {
		lo := LobbyOption{[]int{2}, v.Buffer()}
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
			l.Games = lo.Option
			for i, game := range lo.Option {
				fmt.Fprintln(v, strconv.Itoa(i+1)+" - Partida "+strconv.Itoa(game))
			}
		}
	}
	return nil
}

func (l *Lobby) waitPlayer() {
	var start LobbyOption
	if err := l.Decoder.Decode(&start); err == nil {
		ga := game.NewGame(l.G, l.Encoder, l.Decoder)
		ga.Run()
	}
}

func (l *Lobby) reconect() {
	l.Conn.Close()
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
