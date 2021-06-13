package game_test

import (
	"fmt"
	"testing"
	"time"

	"tp.app/server/game"
	"tp.app/server/user"
)

/*
	Struct y metodos auxiliares
*/

func user_recv(u *user.User) {
	for {
		fmt.Printf("%d: %s\n", u.Id, <-u.Send_chan)
	}
}

/*
	Tests
*/

func TestEchoDeMensaje(t *testing.T) {
	test_users_chan := make(chan user.User)

	san := user.User{Id: 1, Send_chan: make(chan string), Recv_chan: make(chan string)}
	ele := user.User{Id: 2, Send_chan: make(chan string), Recv_chan: make(chan string)}
	ari := user.User{Id: 3, Send_chan: make(chan string), Recv_chan: make(chan string)}

	go func() {
		test_users_chan <- san
		test_users_chan <- ele
		test_users_chan <- ari
	}()

	go func() {
		go func() { san.Recv_chan <- "San" }()
		go user_recv(&san)
	}()
	go func() {
		go func() { ele.Recv_chan <- "Ele" }()
		go user_recv(&ele)
	}()
	go func() {
		go func() { ari.Recv_chan <- "exit" }()
		go user_recv(&ari)
	}()

	game.Start(&test_users_chan)
	time.Sleep(time.Second)
}
