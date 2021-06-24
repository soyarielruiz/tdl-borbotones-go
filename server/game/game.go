package game

import (
	"log"
	"math/rand"
	"time"

	"github.com/soyarielruiz/tdl-borbotones-go/server/action"
	"github.com/soyarielruiz/tdl-borbotones-go/server/stack"
	"github.com/soyarielruiz/tdl-borbotones-go/server/user"
)

type Card struct {
	Number int
	Suit   string
}

type Game struct {
	UserChan    chan user.User
	Users       map[string]user.User
	Deck        *stack.Stack
	DiscardPile *stack.Stack
}

func (g *Game) Init() {
	g.Deck = stack.New()
	g.DiscardPile = stack.New()
	numbers := [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	suits := [4]string{"red", "green", "blue", "yellow"}
	for _, s := range suits {
		for _, n := range numbers {
			g.DiscardPile.Push(Card{n, s})
		}
	}
	g.shuffle()
}

func (g *Game) shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(g.DiscardPile.Size(), func(i, j int) { g.DiscardPile.Swap(i, j) })
	g.Deck.Push(g.DiscardPile)
	g.DiscardPile.Clear()
}

func Start(uc chan user.User, gameNumber int) {
	log.Printf("Initializing game number: %d\n", gameNumber)
	g := Game{UserChan: uc, Users: make(map[string]user.User)}
	g.recvUsers(gameNumber)
	exit := false
	for !exit {
		for k, v := range g.Users {
			log.Printf("Waiting for action from %s in game %d\n", k, gameNumber)
			actionToApply := g.recvMsg(&v)
			g.sendMsg(&actionToApply)
			if actionToApply.Command == "exit" {
				log.Printf("Exit command received in game %d\n", gameNumber)
				exit = true
				return
			}
		}
	}
}

func (g *Game) recvUsers(number int) {
	for {
		u := <-g.UserChan
		log.Printf("New usr received in game %d. %s", number, u)
		g.Users[u.PlayerId] = u
		if len(g.Users) == 3 {
			log.Printf("3 users connect to game %d. Starting game", number)
			return
		} else {
			log.Printf("No enough users connected to game %d for start the game", number)
		}
	}
}

func (g *Game) sendMsg(a *action.Action) {
	for _, u := range g.Users {
		u.SendChannel <- *a
	}
}

func (g *Game) recvMsg(u *user.User) action.Action {
	return <-u.ReceiveChannel
}
