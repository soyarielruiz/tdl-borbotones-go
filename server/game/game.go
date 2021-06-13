package game

import "tp.app/server/user"

type Game struct {
	User_chan *chan user.User
	Users     []user.User
}

func Start(uc *chan user.User) {
	g := Game{User_chan: uc, Users: make([]user.User, 0)}
	RecvUsers(&g)
	i := 0
	for {
		msg := RecvMsg(&g, &g.Users[i])
		SendMsg(&g, &msg)
		if msg == "exit" {
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
		g.Users = append(g.Users, <-*g.User_chan)
		if len(g.Users) == 3 {
			return
		}
	}
}

func SendMsg(g *Game, s *string) {
	for _, u := range g.Users {
		u.Send_chan <- *s
	}
}

func RecvMsg(g *Game, u *user.User) string {
	return <-u.Recv_chan
}
