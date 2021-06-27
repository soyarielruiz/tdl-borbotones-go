package game

import (
	"github.com/soyarielruiz/tdl-borbotones-go/server/deck"
	"github.com/soyarielruiz/tdl-borbotones-go/tools"
	"log"
	"fmt"

	"github.com/soyarielruiz/tdl-borbotones-go/server/user"
)

type Game struct {
	UserChan           <-chan user.User
	Users              map[string]user.User
	Deck               deck.Deck
	RecvChan           chan tools.Action
	CommandHandler     map[tools.Command]CommandHandler
	Ended              bool
	Started            bool
	ActualUserIdToPlay string
	GameNumber         int
}

func NewGame(userChannel chan user.User, gameNumber int) *Game {
	game := Game{UserChan: userChannel, Users: make(map[string]user.User), RecvChan: make(chan tools.Action)}
	game.GameNumber = gameNumber
	game.Deck = *deck.NewDeck()
	game.Ended = false
	game.Started = false
	game.CommandHandler = make (map[tools.Command]CommandHandler)
	game.CommandHandler[tools.DROP] = DropHandler{}
	game.CommandHandler[tools.EXIT] = ExitHandler{}
	game.CommandHandler[tools.TAKE] = TakeHandler{}
	return &game
}

func (game *Game) Run() {
	log.Printf("Initializing game number: %d\n", game.GameNumber)
	game.recvUsers()
	game.Started = true
	var start tools.Action
	game.sendToAll(&start)
	game.sendInitialCards()
	for !game.Ended {
		action := <-game.RecvChan
		fmt.Println(action)
		game.CommandHandler[action.Command].Handle(action, game)
	}
	game.closeAll()
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
			log.Printf("No enough users connected to game %d for start the game", game.GameNumber)
		}
	}
}

func (game *Game) sendToAll(a *tools.Action) {
	for _, u := range game.Users {
		u.SendChannel <- *a
	}
}

func (game *Game) closeAll() {
	log.Printf("Close All in game %d\n", game.GameNumber)
	for _, u := range game.Users {
		close(u.SendChannel)
	}
	close(game.RecvChan)
}

func (game *Game) IsUserTurn(id string) bool {
	return game.ActualUserIdToPlay == id
}

func (game *Game) sendInitialCards() {
	for _, u := range game.Users {
		u.SendChannel <- tools.Action{"", tools.Card{}, u.PlayerId, "", game.Deck.GetCardsFromDeck(3)}
	}
}
