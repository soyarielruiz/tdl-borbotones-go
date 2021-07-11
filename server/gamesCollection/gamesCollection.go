package gamesCollection

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/soyarielruiz/tdl-borbotones-go/server/game"
	"github.com/soyarielruiz/tdl-borbotones-go/server/user"
)

type GamesCollection struct {
	games         map[int]*game.Game
	gameNumber    int
	gamesChannels map[int]chan *user.User
	mu            sync.Mutex
}

type LobbyOption struct {
	Option []int `json:"option"`
}

func NewCollection() *GamesCollection {
	return &GamesCollection{gameNumber: 0, gamesChannels: make(map[int]chan *user.User),
		games: make(map[int]*game.Game)}
}

func (collection *GamesCollection) CreateNewGame(conn net.Conn, nick string) {
	fmt.Println("entre a crear nuevo juego")
	collection.gameNumber = collection.gameNumber + 1
	users := make(chan *user.User)
	new_game := game.NewGame(users, collection.gameNumber)
	collection.games[collection.gameNumber] = new_game
	collection.gamesChannels[collection.gameNumber] = users
	go new_game.Run()
	users <- user.NewUser(conn, nick)
}

func (collection *GamesCollection) SendExistingGames(conn net.Conn, encoder *json.Encoder) int {
	collection.mu.Lock()
	var games []int
	for game_id, game := range collection.games {
		if !game.Started {
			games = append(games, game_id)
		}
	}
	gameOption := LobbyOption{games}
	encoder.Encode(&gameOption)
	defer collection.mu.Unlock()
	return len(games)
}

func (collection GamesCollection) AddUserToGame(conn net.Conn, gameId int, nick string) bool {
	collection.mu.Lock()
	if collection.games[gameId].IsAvailableToJoin() {
		gameChannel := collection.gamesChannels[gameId]
		gameChannel <- user.NewUser(conn, nick)
		defer collection.mu.Unlock()
		return true
	}
	defer collection.mu.Unlock()
	return false
}

func (collection GamesCollection) DeleteDeadGames() {
	collection.mu.Lock()
	for game_id, game := range collection.games {
		if game.Ended {
			delete(collection.games, game_id)
			delete(collection.gamesChannels, game_id)
		}
	}
	collection.mu.Unlock()
}

func (collection GamesCollection) AreAllGamesFinished() <-chan bool {
	gamesFinished := make(chan bool)
	go func() {
		defer close(gamesFinished)
		for {
			finishedGames := 0
			for _, game := range collection.games {
				if game.Ended {
					finishedGames = finishedGames + 1
				}
			}
			if finishedGames == len(collection.games) {
				break
			}
			time.Sleep(5 * time.Second)
		}
		gamesFinished <- true
	}()
	return gamesFinished
}
