package game

import (
	"fmt"

	"github.com/soyarielruiz/tdl-borbotones-go/server/action"
	"github.com/soyarielruiz/tdl-borbotones-go/server/user"
)

type Game struct {
	User_chan chan user.User
	Users     []user.User
}

func Start(uc chan user.User, game_number int) {
	fmt.Printf("Game number: %d\n", game_number)
	g := Game{User_chan: uc, Users: make([]user.User, 0)}
	RecvUsers(&g)
	i := 0
	for {
		fmt.Printf("Turno de %d\n", i+1)
		action := RecvMsg(&g, &g.Users[i])
		SendMsg(&g, &action)
		if action.Command == "exit" {
			return
		}
		i++
		if i == 3 {
			i = 0
		}
	}
}

func RecvUsers(g *Game) {
	for {
		u := <-g.User_chan
		fmt.Println(u.PlayerId)
		g.Users = append(g.Users, u)
		if len(g.Users) == 3 {
			return
		}
	}
}

func SendMsg(g *Game, a *action.Action) {
	for _, u := range g.Users {
		u.SendChannel <- *a
	}
}

func RecvMsg(g *Game, u *user.User) action.Action {
	return <-u.ReceiveChannel
}
