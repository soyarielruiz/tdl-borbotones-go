package gameManager

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/soyarielruiz/tdl-borbotones-go/server/gamesCollection"
)

const (
	connHost = "localhost"
	connPort = "8080"
	connType = "tcp"
)

type LobbyOption struct{
	Option []int `json:"option"`
	Nickname string `json:"nickname"`
}

type GameManager struct {
	listener net.Listener
	collection *gamesCollection.GamesCollection 
}

func NewGameManager() (*GameManager) {
	 return &GameManager{collection:gamesCollection.NewCollection()}
}

func (manager *GameManager) Start() {
	log.Println("Starting " + connType + " server on " + connHost + ":" + connPort)

	server, err := net.Listen(connType, connHost+":"+connPort)
	if err != nil {
		log.Printf("Error listening:", err.Error())
		os.Exit(1)
	}
	manager.listener = server
	defer server.Close()

	manager.acceptConnections()
}

func (manager *GameManager) Stop() {
	manager.listener.Close()
	log.Printf("Waiting for games to finish...\n")
	games_finished := <-manager.collection.AreAllGamesFinished()
	if games_finished {
		log.Printf("All games finished. Quitting server \n")
	}
}

func (manager *GameManager) acceptConnections() {
	for {
		client, err := manager.listener.Accept()
		if err != nil {
			log.Printf("Error connecting:", err.Error())
			break
		}
		log.Printf("New connection accepted from %s\n", client.RemoteAddr())
		go lobby(client,manager.collection)
	}
}

func lobby(conn net.Conn, collection *gamesCollection.GamesCollection) {
	 decoder := json.NewDecoder(conn)
	 var gameOption LobbyOption
	 for{
		if error:=decoder.Decode(&gameOption); error==nil{
			option:=gameOption.Option[0]
			switch option {
			case 1 : 
				collection.CreateNewGame(conn, gameOption.Nickname)
			case 2: 
				success:=joinExistingGame(conn,collection, gameOption.Nickname)
				if success {
					break 
				} else {
					continue 
				}
			}
		} else {
			conn.Close()
		}
		break
	}
	fmt.Println("sali del lobby")
}

func joinExistingGame(conn net.Conn,collection *gamesCollection.GamesCollection, nick string) bool {
	collection.DeleteDeadGames()
	games:=collection.SendExistingGames(conn)
	if games !=0 {
		decoder := json.NewDecoder(conn)
		var gameNumber LobbyOption
		if error:=decoder.Decode(&gameNumber); error!=nil {
			return false
		}
		collection.AddUserToGame(conn,gameNumber.Option[0], nick)
		return true
	} else {
		return false
	}
}