package gamesCollection

import (
	"net"
	"encoding/json"
	"github.com/soyarielruiz/tdl-borbotones-go/server/game"
	"github.com/soyarielruiz/tdl-borbotones-go/server/user"
)

type GamesCollection struct{
	Games map[int] *game.Game
	Game_number int
	GamesChannels map[int] chan user.User
}

type LobbyOption struct{
	Option []int `json:"option"`
}

func NewCollection() *GamesCollection{
	 return &GamesCollection{Game_number:0,GamesChannels: make(map[int] chan user.User),
		                     Games:make(map[int] *game.Game)}
}

func (collection *GamesCollection) CreateNewGame(conn net.Conn){
	  collection.Game_number=collection.Game_number+1
	  users:=make(chan user.User)
	  new_game:=game.NewGame(users,collection.Game_number)
	  collection.Games[collection.Game_number]=new_game
	  collection.GamesChannels[collection.Game_number]=users
	  go new_game.Run()
	  users <- user.CreateFromConnection(conn)
}

func (collection *GamesCollection) SendExistingGames(conn net.Conn){
	var games []int
	for game_id,game:= range collection.Games{
		if(!game.Started){
			games=append(games,game_id)
		}
	}
	encoder := json.NewEncoder(conn)
	gameOption:=LobbyOption{games}
	encoder.Encode(&gameOption)
}

func (collection GamesCollection) AddUserToGame(conn net.Conn, gameId int){
	gameChannel:=collection.GamesChannels[gameId]
	gameChannel <- user.CreateFromConnection(conn)
}

func (collection GamesCollection) DeleteDeadGames(){
	for game_id,game := range collection.Games{
		if game.Ended {
			delete(collection.Games,game_id)
			delete(collection.GamesChannels,game_id)
		}
	}
}

