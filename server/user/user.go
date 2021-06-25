package user

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/soyarielruiz/tdl-borbotones-go/tools"

	"github.com/google/uuid"
)

type User struct {
	SendChannel    chan tools.Action
	ReceiveChannel chan tools.Action
	PlayerId       string
	conn           net.Conn
}

func (u User) String() string {
	return fmt.Sprintf("PlayerId:\"%s\", conn:\"%s\"", u.PlayerId, u.conn.RemoteAddr())
}

func CreateFromConnection(conn net.Conn) User {
	var usr = User{make(chan tools.Action), make(chan tools.Action), uuid.New().String(), conn}
	go Send(usr)
	// go Receive(usr)
	return usr
}

func Send(usr User) {
	encoder := json.NewEncoder(usr.conn)
	for {
		action, ok := <-usr.SendChannel
		if !ok {
			log.Println("Send Close chan")
			break
		}
		// log.Printf("Sending action to usr %s. %s", usr, action)
		encoder.Encode(&action)
	}
}

func Receive(usr User) {
	decoder := json.NewDecoder(usr.conn)
	for {
		var action tools.Action
		decoder.Decode(&action)
		// log.Printf("Receive action from usr %s.\n %s", usr, action)
		usr.ReceiveChannel <- action
	}
}
