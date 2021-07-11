package game

import (
	"log"
	"time"

	"github.com/soyarielruiz/tdl-borbotones-go/server/deck"
	"github.com/soyarielruiz/tdl-borbotones-go/server/turnero"
	"github.com/soyarielruiz/tdl-borbotones-go/tools"

	"github.com/soyarielruiz/tdl-borbotones-go/server/user"
)

type Game struct {
	UserChan       <-chan *user.User
	Users          map[string]*user.User
	Deck           deck.Deck
	RecvChan       chan tools.Action
	CommandHandler map[tools.Command]CommandHandler
	Ended          bool
	Started        bool
	GameNumber     int
	Tur            turnero.Turnero
}

func NewGame(userChannel chan *user.User, gameNumber int) *Game {
	game := Game{UserChan: userChannel, Users: make(map[string]*user.User), RecvChan: make(chan tools.Action)}
	game.GameNumber = gameNumber
	game.Deck = *deck.NewDeck()
	game.Ended = false
	game.Started = false
	game.CommandHandler = make(map[tools.Command]CommandHandler)
	game.CommandHandler[tools.DROP] = DropHandler{}
	game.CommandHandler[tools.EXIT] = ExitHandler{}
	game.CommandHandler[tools.TAKE] = TakeHandler{}
	return &game
}

func (game *Game) Run() {
	log.Printf("Initializing game number: %d\n", game.GameNumber)
	game.recvUsers()
	if !game.Ended {
		game.Started = true
		game.Tur = *turnero.New(game.Users)
		var start tools.Action
		game.SendToAll(&start)
		game.sendInitialCards()
		for !game.Ended {
			action := <-game.RecvChan
			if action.Command != "" {
				game.CommandHandler[action.Command].Handle(action, game)
			}
		}
		log.Printf("Game %d ended", game.GameNumber)
	}
}

func (game *Game) recvUsers() {
	for {
		select {
			case u := <-game.UserChan:
				game.addUser(u)
			case <-time.After(10 * time.Second):
				game.closeAll()
				game.Ended = true
				return
		}
		if len(game.Users) == 3 {
			log.Printf("3 users connect to game %d. Starting game", game.GameNumber)
			return
		} else {
			log.Printf("No enough users connected to game %d for start the game", game.GameNumber)
		}
	}
}

func (game *Game) addUser(u *user.User){
	u.ReceiveChannel = game.RecvChan
	go u.Receive()
	log.Printf("New usr received in game %d. %s", game.GameNumber, u)
	game.Users[u.PlayerId] = u
}

func (game *Game) SendToAll(a *tools.Action) {
	for _, u := range game.Users {
		u.SendChannel <- *a
	}
}

func (game *Game) closeAll() {
	log.Printf("Close All in game %d\n", game.GameNumber)
	for _, u := range game.Users {
		u.Close()
	}
	close(game.RecvChan)
}

func (game *Game) TurnMoveForward() {
	game.Tur.Next()
	game.Users[game.Tur.CurrentUser()].SendChannel <- tools.Action{
		Command:  tools.TURN_ASSIGNED,
		Card:     game.Deck.GetFrontCard(),
		PlayerId: game.Tur.CurrentUser(),
		Message:  "It's your turn to play!",
		Cards:    nil,
	}
}

func (game *Game) sendInitialCards() {
	for _, u := range game.Users {
		cardsAction := tools.Action{Card: game.Deck.GetFrontCard(), PlayerId: u.PlayerId, Cards: game.Deck.GetCardsFromDeck(3)}
		u.SendChannel <- cardsAction
	}
	game.Users[game.Tur.CurrentUser()].SendChannel <- tools.Action{
		Command:  tools.TURN_ASSIGNED,
		Card:     game.Deck.GetFrontCard(),
		PlayerId: game.Tur.CurrentUser(),
		Message:  "It's your turn to play!",
		Cards:    nil,
	}
}

func (game *Game) IsAvailableToJoin() bool {
	 return (!game.Ended && !game.Started)
}
