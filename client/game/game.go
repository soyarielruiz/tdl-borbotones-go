package game

import (
	"encoding/json"
	"net"
	"log"
	"fmt"
	"os"
)

const (
	serverAddress = "127.0.0.1"
	serverPort    = "8080"
	serverConn    = "tcp"
)

type Game struct{
	Encoder *json.Encoder
	Decoder *json.Decoder
	Conn *net.TCPConn
}

func NewGame() *Game{
	 conn:=startClient()
	 encoder:=json.NewEncoder(conn)
	 decoder:=json.NewDecoder(conn)
	 var game=Game{encoder,decoder,conn}
	 return &game
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