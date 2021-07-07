package game

import (
	"encoding/json"
	"log"
	"net"
)

const (
	serverAddress = "127.0.0.1"
	serverPort    = "8080"
	serverConn    = "tcp"
)

type Game struct {
	Encoder *json.Encoder
	Decoder *json.Decoder
	Conn    *net.TCPConn
	// G       *gocui.Gui
}

func NewGame() *Game {
	conn := startClient()
	encoder := json.NewEncoder(conn)
	decoder := json.NewDecoder(conn)
	var game = Game{encoder, decoder, conn}
	return &game
}

func startClient() *net.TCPConn {

	serverConnection := serverAddress + ":" + serverPort
	log.Println("Starting " + serverConn + " client on " + serverConnection)
	tcpAddr, err := net.ResolveTCPAddr(serverConn, serverAddress+":"+serverPort)
	if err != nil {
		log.Panic(err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Panic(err)
	}

	err = conn.SetWriteBuffer(10)
	if err != nil {
		log.Panic(err)
	}

	return conn
}
