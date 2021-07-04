package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/awesome-gocui/gocui"
	"github.com/soyarielruiz/tdl-borbotones-go/client/translator"
	"github.com/soyarielruiz/tdl-borbotones-go/tools"
	"log"
	"net"
	"os"
	"strconv"
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
			messageToSend, err := translator.CreateAnAction(messageToUse)
			if err != nil {
				out, _ := g.View("mesa")
				fmt.Fprintf(out, "Error al crear la accion, probar nuevamente")
			}

			err = encoder.Encode(&messageToSend)
			checkError(err)

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
	})

	if err != nil {
		log.Println("Cannot bind the enter key:", err)
	}

	// receiving from server
	go receivingData(g, conn)

	if err := view.InitKeybindings(g); err != nil {
		log.Fatalln(err)
	}

	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Panicln(err)
	}
}

func receivingData(g *gocui.Gui, conn *net.TCPConn) {
	for {
		time.Sleep(1*time.Second)
		if err := view.ReceiveMsgFromGame(g, conn); err != nil {
			log.Fatalln(err)
		}
	}
}

func lobby(conn *net.TCPConn ){
	fmt.Println("* * * * * * * * * *")
	fmt.Println("* WELCOME TO GUNO *")
	fmt.Println("* * * * * * * * * *")
	encoder := json.NewEncoder(conn)
	decoder := json.NewDecoder(conn)
	for {
		input:=initialOption()
		option :=LobbyOption{[]int{input}}
		encoder.Encode(&option)
		if input==2{
			var games LobbyOption
			decoder.Decode(&games)
			if (len(games.Option)==0){
				fmt.Println("There are no current games available.")
				continue
			}
			input=gameOption(games)
			option2 :=LobbyOption{[]int{input}}
			encoder.Encode(&option2)
		}
		break
	}
	fmt.Println("Waiting for new members to start")
	var start tools.Action
	decoder.Decode(&start)
}

func initialOption() int {
	var option int
	var s string
	for {
		fmt.Println("1:Start new game\n2:Join game")
		_, err := fmt.Scan(&s)
        option, err = strconv.Atoi(s)
		if option!=1 && option!=2 || err!=nil {
			fmt.Println("Wrong option. Please type 1 or 2: ")
		}else{
			break
		}
	}
	return option
}

func gameOption(games LobbyOption) int {
	var option int
	var s string
	for {
		fmt.Println("Choose game number:")
		fmt.Println(games.Option)
		_, err := fmt.Scan(&s)
        option, err = strconv.Atoi(s)
		if err!=nil || !checkGameId(games,option) {
			fmt.Println("Wrong option")
		}else{
			break
		}
	}
	return option
}

func checkGameId(games LobbyOption,option int) bool{
	for _, gameId := range games.Option {
        if gameId == option {
            return true
        }
    }
    return false
}

func startClient() *net.TCPConn {

	serverConnection := serverAddress + ":" + serverPort

	log.Println("Starting " + serverConn + " client on " + serverConnection)

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
