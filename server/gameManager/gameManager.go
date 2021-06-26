package gameManager

import (
	"fmt"
	"log"
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
	log.Println("Starting " + connType + " server on " + connHost + ":" + connPort)

	server, err := net.Listen(connType, connHost+":"+connPort)
	if err != nil {
		log.Fatalln("Error listening:", err.Error())
		os.Exit(1)
	}

	defer server.Close()
	acceptConnections(server)
}

func acceptConnections(listener net.Listener) {
	user_counter := 0
	game_number := 1
	users := make(chan user.User)
	go game.Run(users, game_number)
	for {
		client, err := listener.Accept()
		log.Printf("New connection accepted from %s\n", client.RemoteAddr())
		if err != nil {
			log.Fatalln("Error connecting:", err.Error())
			return
		}
		if user_counter == 3 {
			log.Printf("New game started %d", game_number)
			game_number = game_number + 1
			users = make(chan user.User)
			go game.Run(users, game_number)
			user_counter = 0
		}
		go handleConnection(client, users)
		user_counter = user_counter + 1
	}
}

func handleConnection(conn net.Conn, users chan user.User) {
	message := []byte("WELCOME TO GUNO\n")
	_, err := conn.Write(message)
	checkError(err)
	users <- user.CreateFromConnection(conn)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}