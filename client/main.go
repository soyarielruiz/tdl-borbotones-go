package main

import (
	"encoding/json"
	"fmt"
	"github.com/awesome-gocui/gocui"
	"github.com/soyarielruiz/tdl-borbotones-go/client/translator"
	"github.com/soyarielruiz/tdl-borbotones-go/tools"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/soyarielruiz/tdl-borbotones-go/client/view"
)

const (
	serverAddress = "127.0.0.1"
	serverPort    = "8080"
	serverConn    = "tcp"
)

type LobbyOption struct{
	Option []int `json:"option"`
}

func main() {
	
	conn := startClient()
	
	lobby(conn)

	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen

	//conn := startClient()

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

	// receiving from server
	go receivingData(g, conn)

	if err := view.InitKeybindings(g); err != nil {
		log.Fatalln(err)
	}

}

func receivingData(g *gocui.Gui, conn *net.TCPConn) {
	if err := ReceiveMsgFromGame(g, conn); err != nil {
		log.Fatalln(err)
	}
}

func lobby(conn *net.TCPConn ){
	fmt.Println("* * * * * * * * * *\n")
	fmt.Println("* WELCOME TO GUNO *\n")
	fmt.Println("* * * * * * * * * *\n")
	fmt.Println("1:Start new game\n2:Join game\n")
	var input int
    fmt.Scanf("%d", &input)
	option :=LobbyOption{[]int{input}}
	encoder := json.NewEncoder(conn)
	encoder.Encode(&option)
	decoder := json.NewDecoder(conn)
	if input==2{
		//recibo las partidas disponibles
		var games LobbyOption
		decoder.Decode(&games)
		fmt.Println("Choose game number:\n")
		fmt.Println(games)
		//mando que me quiero conectar a la partida
		fmt.Scanf("%d", &input)
		option2 :=LobbyOption{[]int{input}}
		encoder.Encode(&option2)
	}
	fmt.Println("Waiting for new members to start")
	var start tools.Action
	decoder.Decode(&start)
}

func startClient() *net.TCPConn {

	serverConnection := serverAddress + ":" + serverPort

	log.Println("Starting " + serverConn + " client on " + serverConnection +"\n")

	tcpAddr, err := net.ResolveTCPAddr(serverConn, serverAddress+":"+serverPort)
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

	return tools.Action{command, card, "", message,[] tools.Card{}}
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

func ReceiveMsgFromGame(g *gocui.Gui, conn *net.TCPConn) error {
	//wait a starting moment
	time.Sleep(1*time.Second)
	for {
		decoder := json.NewDecoder(conn)
		var action tools.Action
		decoder.Decode(&action)

		if len(action.Command.String()) > 1 {
			out, _ := g.View("mesa")
			//fmt.Fprintln(out, "message2 : ", action.Command.String())
			message, err := translator.TranslateMessageFromServer(action)
			if err == nil {
				fmt.Fprintln(out, message)
				message = ""
			}
		}
	}
}
