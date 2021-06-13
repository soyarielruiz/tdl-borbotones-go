package gameManager

import (
	"fmt"
	"net"
	"os"

	"github.com/soyarielruiz/tdl-borbotones-go/server/game"
	"github.com/soyarielruiz/tdl-borbotones-go/server/user"
)

const (
	connHost = "localhost"
	connPort = "8080"
	connType = "tcp"
)

func Start() {
	fmt.Println("Starting " + connType + " server on " + connHost + ":" + connPort)

	server, err := net.Listen(connType, connHost+":"+connPort)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	defer server.Close()
	acceptConnections(server)
}

func acceptConnections(listener net.Listener) {
	user_counter := 0
	game_number := 1
	users := make(chan user.User)
	go game.Start(users, game_number)
	for {
		client, err := listener.Accept()
		if err != nil {
			fmt.Println("Error connecting:", err.Error())
			return
		}
		if user_counter == 3 {
			game_number = game_number + 1
			users = make(chan user.User)
			go game.Start(users, game_number)
			user_counter = 0
		}
		go handleConnection(client, users)
		user_counter = user_counter + 1
	}
}

func handleConnection(conn net.Conn, users chan user.User) {
	message := []byte("WELCOME TO GUNO\n")
	conn.Write(message)
	users <- user.CreateFromConnection(conn)
}
