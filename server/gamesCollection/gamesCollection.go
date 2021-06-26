package gamesCollection

import (
	"net"
	"encoding/json"
	"github.com/soyarielruiz/tdl-borbotones-go/server/game"
	"github.com/soyarielruiz/tdl-borbotones-go/server/user"
	"fmt"
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

func (collection *GamesCollection) CreateDummyGames(){
	 users1:=make(chan user.User)
	 users2:=make(chan user.User)
	 game1:=game.NewGame(users1,1)
	 game2:=game.NewGame(users2,2)
	 collection.GamesChannels[1]=users1
	 collection.GamesChannels[2]=users2
	 collection.Games[1]=game1
	 collection.Games[2]=game2
	 go game1.Run()
	 go game2.Run()
}

func (collection *GamesCollection) CreateNewGame(conn net.Conn){
	  collection.Game_number=collection.Game_number+1
	  fmt.Println("creamos nueva partida numero: ", collection.Game_number)
	  users:=make(chan user.User)
	  new_game:=game.NewGame(users,collection.Game_number)
	  collection.Games[collection.Game_number]=new_game
	  collection.GamesChannels[collection.Game_number]=users
	  go new_game.Run()
	  users <- user.CreateFromConnection(conn)
}

func (collection *GamesCollection) SendExistingGames(conn net.Conn){
	fmt.Println("sending available games")
	var games []int
	for game_id,_ := range collection.Games{
		games=append(games,game_id)
	}
	fmt.Println(games)
	encoder := json.NewEncoder(conn)
	gameOption:=LobbyOption{games}
	encoder.Encode(&gameOption)
}

func (collection GamesCollection) AddUserToGame(conn net.Conn, gameId int){
	gameChannel:=collection.GamesChannels[gameId]
	gameChannel <- user.CreateFromConnection(conn)
}

func (collection GamesCollection) deleteGames(){}

