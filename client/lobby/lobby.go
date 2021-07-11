package lobby

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	"github.com/awesome-gocui/gocui"
	"github.com/soyarielruiz/tdl-borbotones-go/client/game"
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

type UserJoined struct {
	Success int `json:"success"`
}

func New(g *gocui.Gui) (*Lobby, error) {
	tcpAddr, err := net.ResolveTCPAddr(serverConn, serverAddress+":"+serverPort)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}
	conn.SetWriteBuffer(10)
	return &Lobby{g, conn, json.NewEncoder(conn), json.NewDecoder(conn), nil}, nil
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

	if err := g.SetKeybinding("lost", gocui.KeyF3, gocui.ModNone, l.Quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("wait", gocui.KeyF3, gocui.ModNone, l.Back); err != nil {
		return err
	}

	if err := g.SetKeybinding("timeout", gocui.KeyF3, gocui.ModNone, l.Back); err != nil {
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
			lo := LobbyOption{[]int{l.Games[key-1]}, ""}
			l.Encoder.Encode(lo)
			var suc UserJoined
			l.Decoder.Decode(&suc)
			g.DeleteView("gamelist")
			if suc.Success == 50 {
				v, _ := g.SetView("wait", 27, 1, 52, 4, gocui.TOP)
				g.SetCurrentView(v.Name())
				v.Title = "Waiting players"
				fmt.Fprintln(v, "F3 - Back")
				fmt.Fprintf(v, "Selected option %d\n", key)
				go l.waitPlayer()
			} else {
				v, _ := g.SetView("wait", 27, 1, 52, 5, gocui.TOP)
				v.Title = "Error"
				g.SetCurrentView(v.Name())
				fmt.Fprintln(v, "F3 - Back")
				fmt.Fprintln(v, "Game selected not")
				fmt.Fprintln(v, "exists")
			}
		}
		return nil
	}
}

func (l *Lobby) Back(g *gocui.Gui, v *gocui.View) error {
	l.reconect()

	if g.CurrentView().Name() == "timeout" {
		if err := g.DeleteView("timeout"); err != nil {
			return err
		}
		if err := g.DeleteView("wait"); err != nil {
			return err
		}
	}

	if g.CurrentView().Name() == "wait" {
		if err := g.DeleteView("wait"); err != nil {
			return err
		}
	}

	if g.CurrentView().Name() == "gamelist" {
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
		v.Title = "Welcome to GUNO"
		fmt.Fprintln(v, "F1 - New game")
		fmt.Fprintln(v, "F2 - Find game")
		fmt.Fprintln(v, "F3 - Exit")
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
		err := l.Encoder.Encode(lo)
		if err != nil {
			v, _ := g.SetView("lost", 1, 9, 25, 11, gocui.TOP)
			g.SetCurrentView(v.Name())
			v.Title = "Server"
			fmt.Fprintln(v, "Connection lost")
		} else {
			v, _ := g.SetView("wait", 27, 1, 52, 4, gocui.TOP)
			g.SetCurrentView(v.Name())
			v.Title = "Waiting players"
			fmt.Fprintln(v, "F3 - Back")
			fmt.Fprintln(v, "New game")
			go l.waitPlayer()
		}
	}
	return nil
}

func (l *Lobby) FindGame(g *gocui.Gui, v *gocui.View) error {
	nick := v.Buffer()
	if nick != "" {
		lo := LobbyOption{[]int{2}, v.Buffer()}
		l.Encoder.Encode(lo)
		err := l.Decoder.Decode(&lo)
		if err != nil {
			v, _ := g.SetView("lost", 1, 9, 25, 11, gocui.TOP)
			g.SetCurrentView(v.Name())
			v.Title = "Server"
			fmt.Fprintln(v, "Connection lost")
		} else {
			_, y := g.Size()
			v, _ := g.SetView("gamelist", 27, 1, 52, y-3, gocui.TOP)
			g.SetCurrentView(v.Name())
			v.Title = "Game list"
			fmt.Fprintln(v, "F3 - Back")
			fmt.Fprintln(v, "-----------------------")
			l.Games = lo.Option
			for i, game := range lo.Option {
				fmt.Fprintln(v, strconv.Itoa(i+1)+" - Game "+strconv.Itoa(game))
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
	} else {
		l.G.Update(l.timeout())
	}
}

func (l *Lobby) timeout() func(g *gocui.Gui) error {
	return func(g *gocui.Gui) error {
		v, _ := g.SetView("timeout", 27, 5, 52, 10, gocui.TOP)
		v.Title = "Timeout"
		g.SetCurrentView(v.Name())
		fmt.Fprintln(v, "F3: Back")
		fmt.Fprintln(v, "Time limit expired")
		fmt.Fprintln(v, "Not enough players")
		fmt.Fprintln(v, "connected")
		return nil
	}
}

func (l *Lobby) reconect() error {
	l.Conn.Close()
	tcpAddr, err := net.ResolveTCPAddr(serverConn, serverAddress+":"+serverPort)
	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}
	conn.SetWriteBuffer(10)
	l.Conn = conn
	l.Encoder = json.NewEncoder(conn)
	l.Decoder = json.NewDecoder(conn)
	return nil
}
