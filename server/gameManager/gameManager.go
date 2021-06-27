package gameManager

import (
	"log"
	"net"
	"os"
	"encoding/json"

	"github.com/soyarielruiz/tdl-borbotones-go/server/gamesCollection"
)

const (
	connHost = "localhost"
	connPort = "8080"
	connType = "tcp"
)

type LobbyOption struct{
	Option []int `json:"option"`
}

func Start() {
	log.Println("Starting " + connType + " server on " + connHost + ":" + connPort)

	server, err := net.Listen(connType, connHost+":"+connPort)
	if err != nil {
		log.Fatalln("Error listening:", err.Error())
		os.Exit(1)
	}

	defer server.Close()
	acceptConnections(server)
}

func acceptConnections(listener net.Listener) {
	collection := gamesCollection.NewCollection()
	for {
		client, err := listener.Accept()
		log.Printf("New connection accepted from %s\n", client.RemoteAddr())
		if err != nil {
			log.Fatalln("Error connecting:", err.Error())
			return
		}
		go lobby (client,collection)
	}
}

func lobby(conn net.Conn, collection *gamesCollection.GamesCollection) {
	 decoder := json.NewDecoder(conn)
	 var gameOption LobbyOption
	 decoder.Decode(&gameOption)
	 option:=gameOption.Option[0]
	 switch option {
	 case 1 : 
	 	collection.CreateNewGame(conn)
	 case 2: 
	 	joinExistingGame(conn,collection)
	 }
}

func joinExistingGame(conn net.Conn,collection *gamesCollection.GamesCollection){
	collection.DeleteDeadGames()
	collection.SendExistingGames(conn)
	decoder := json.NewDecoder(conn)
	var gameNumber LobbyOption
	decoder.Decode(&gameNumber)
	collection.AddUserToGame(conn,gameNumber.Option[0])
}