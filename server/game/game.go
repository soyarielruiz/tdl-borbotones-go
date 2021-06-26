package game

import (
	"github.com/soyarielruiz/tdl-borbotones-go/server/commandHandler"
	"github.com/soyarielruiz/tdl-borbotones-go/server/deck"
	"github.com/soyarielruiz/tdl-borbotones-go/server/discardPile"
	"github.com/soyarielruiz/tdl-borbotones-go/tools"
	"log"

	"github.com/soyarielruiz/tdl-borbotones-go/server/user"
)

type Game struct {
	UserChan         <-chan user.User
	Users            map[string]user.User
	Deck             deck.Deck
	DiscardPile      discardPile.DiscardPile
	RecvChan         chan tools.Action
	CommandHandler   map[tools.Command]commandHandler.CommandHandler
	Ended            bool
	NextUserIdToPlay string
	GameNumber       int
}

func (game *Game) Init() {
	game.Deck = deck.Deck{}
	game.Deck.Init()
	game.DiscardPile = discardPile.DiscardPile{}
	game.DiscardPile.Init()
	game.Ended = false
	game.CommandHandler[tools.DROP] = commandHandler.DropHandler{}
	game.CommandHandler[tools.EXIT] = commandHandler.ExitHandler{}
	game.CommandHandler[tools.TAKE] = commandHandler.TakeHandler{}
}

func Run(userChannel chan user.User, gameNumber int) {
	log.Printf("Initializing game number: %d\n", gameNumber)
	game := createGame(userChannel, gameNumber)
	game.recvUsers()
	game.sendInitialCards()
	for !game.Ended {
		action := <-game.RecvChan
		game.CommandHandler[action.Command].Handle(action, &game)
	}
	game.closeAll(gameNumber)
}

func createGame(userChannel chan user.User, gameNumber int) Game {
	game := Game{UserChan: userChannel, Users: make(map[string]user.User), RecvChan: make(chan tools.Action)}
	game.GameNumber = gameNumber
	game.Init()
	return game
}

func (game *Game) recvUsers() {
	for {
		u := <-game.UserChan
		u.ReceiveChannel = game.RecvChan
		go user.Receive(u)
		log.Printf("New usr received in game %d. %s", game.GameNumber, u)
		game.Users[u.PlayerId] = u
		if len(game.Users) == 3 {
			log.Printf("3 users connect to game %d. Starting game", game.GameNumber)
			return
		} else {
			log.Printf("No enough users connected to game %d for start the game", gameNumber)
		}
	}
}

func (game *Game) sendToAll(a *tools.Action) {
	for _, u := range game.Users {
		u.SendChannel <- *a
	}
}

func (game *Game) closeAll(gn int) {
	log.Printf("Close All in game %d\n", gn)
	for _, u := range game.Users {
		close(u.SendChannel)
	}
	close(game.RecvChan)
}

func (game *Game) IsUserTurn(id string) bool {
	return game.NextUserIdToPlay == id
}

func (game *Game) sendInitialCards() {
	for _, u := range game.Users {
		u.SendChannel <- tools.Action{"", }
	}
}
