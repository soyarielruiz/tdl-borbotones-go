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

type UserJoined struct{
	success int `json:"success"`
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
			log.Printf("Closed connection:", err.Error())
			break
		}
		log.Printf("New connection accepted from %s\n", client.RemoteAddr())
		go manager.lobby(client)
	}
}

func (manager *GameManager) lobby(conn net.Conn) {
	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)
	var gameOption LobbyOption
	for {
		if error:=decoder.Decode(&gameOption); error==nil {
			fmt.Fprintf(os.Stderr,"option: %s\n", gameOption)
			if success := manager.executeOption(conn,gameOption,decoder,encoder); success {
				break 
			}
		} else {
			conn.Close()
			break 
		}
	}
	fmt.Println("sali del lobby")
}


func (manager *GameManager) executeOption(conn net.Conn, gameOption LobbyOption, decoder *json.Decoder, encoder *json.Encoder) bool {
	option := gameOption.Option[0]
	nick := gameOption.Nickname
	switch option {
		case 1 : 
			manager.collection.CreateNewGame(conn,nick)
			return true
		case 2: 
			success := manager.joinExistingGame(conn,nick,decoder,encoder)
			return success
		default:
			return false
	}
}

func (manager *GameManager) joinExistingGame(conn net.Conn, nick string, decoder *json.Decoder, encoder *json.Encoder) bool {
	manager.collection.DeleteDeadGames()
	games:=manager.collection.SendExistingGames(conn, encoder)
	if games !=0 {
		var gameNumber LobbyOption
		if error := decoder.Decode(&gameNumber); error!=nil {
			return false
		}
		success := manager.collection.AddUserToGame(conn,gameNumber.Option[0], nick)
		sendIfJoined(success,decoder,encoder)
		return success
	} else {
		return false
	}
}

func sendIfJoined(success bool, decoder *json.Decoder, encoder *json.Encoder) {
	if success {
		status := UserJoined{50}
		//encoder.Encode(&status)
		fmt.Println("status: ", status)
	} else {
		status := UserJoined{90}
		//encoder.Encode(&status)
		fmt.Println("status: ", status)
	}
}