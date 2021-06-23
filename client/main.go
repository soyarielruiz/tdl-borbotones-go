package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/awesome-gocui/gocui"
	tools "github.com/soyarielruiz/tdl-borbotones-go/tools/action"
	"io"
	"strconv"

	//"github.com/soyarielruiz/tdl-borbotones-go/client/translator"
	"strings"

	"github.com/soyarielruiz/tdl-borbotones-go/client/view"
	"log"
	"net"
	"os"
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

			encoder := json.NewEncoder(conn)
			messageToUse := string(iv.Buffer())
			messageToSend := createAnAction(messageToUse)

			err := encoder.Encode(&messageToSend)
			checkError(err)

			//get cursor position
			x, y := iv.Cursor()

			//adding a visual enter
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

	v, err := g.View("mesa")

	go receivingData(g, v, conn)
}

func receivingData(g *gocui.Gui, v *gocui.View, conn *net.TCPConn) {
	if err := ReceiveMsgFromGame(g, v, conn); err != nil {
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

func createAnAction(messageToSend string) tools.Action {
	words := strings.Fields(messageToSend)

	command := getCommandFromMessage(words[0])
	card := getCardFromMessage(words[1], words[2])
	message := strings.Join(words[3:], " ")

	return tools.Action{command, card, "", message}
}

func getCardFromMessage(color string, number string) tools.Card {
	var colorToUse = tools.GREEN
	switch strings.ToLower(color) {
	case string(tools.GREEN):
		colorToUse = tools.GREEN
	case string(tools.YELLOW):
		colorToUse = tools.YELLOW
	case string(tools.RED):
		colorToUse = tools.RED
	case string(tools.BLUE):
		colorToUse = tools.BLUE
	default:
		colorToUse = tools.GREEN
	}

	value, err := strconv.ParseInt(number, 10, 64)
	checkError(err)

	return tools.Card{int(value), colorToUse}
}

func getCommandFromMessage(message string) tools.Command {

	switch strings.ToLower(message) {
	case string(tools.DROP):
		return tools.DROP
	case string(tools.EXIT):
		return tools.EXIT
	case string(tools.TAKE):
		return tools.TAKE
	default:
		return tools.DROP
	}
}

func ReceiveMsgFromGame(g *gocui.Gui, v *gocui.View, conn *net.TCPConn) error {
	// receives msg from server
	var messageFromServer bytes.Buffer
	io.Copy(&messageFromServer, conn)
	fmt.Println("tiene ", messageFromServer)
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

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}