package main

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
	"log"
	"net"
	"os"
	"tp.app/client/translator"
	"tp.app/client/view"
)

const (
	serverAddress = "127.0.0.1"
	serverPort = "8080"
	serverConn = "tcp"
)

func main() {

	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen

	conn := startClient()

	g.SetManagerFunc(view.Layout)

	// Bind enter key to input to send new messages.
	err = g.SetKeybinding("jugador", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, iv *gocui.View) error {
		iv.Autoscroll = true
		iv.Editable = true
		// Read buffer from the beginning.
		iv.Rewind()

		// Send message if text was entered.
		if len(iv.Buffer()) >= 2 {
			msg := translator.ToJSON(iv.Buffer())
			_, err := conn.Write([]byte(msg))
			checkError(err)

			//get cursor position
			x, y := iv.Cursor()

			//adding a enter
			y = y + 1
			x = 0

			err = iv.SetCursorUnrestricted(x, y)
			if err != nil {
				log.Println("Failed to set cursor:", err)
			}
			return err
		}
		return nil
	})

	if err != nil {
		log.Println("Cannot bind the enter key:", err)
	}

	if err := view.InitKeybindings(g); err != nil {
		log.Fatalln(err)
	}
}

func startClient() *net.TCPConn {

	serverConnection := serverAddress + ":" + serverPort

	log.Println("Starting " + serverConn + " client on " + serverConnection)

	tcpAddr, err := net.ResolveTCPAddr(serverConn, serverAddress + ":" + serverPort)
	checkError(err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	err = conn.SetWriteBuffer(10)
	checkError(err)

	return conn
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}