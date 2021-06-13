package user

import (
	"encoding/json"
	"github.com/soyarielruiz/tdl-borbotones-go/server/action"
	"net"

	"github.com/google/uuid"
)

type User struct {
	SendChannel    chan action.Action
	ReceiveChannel chan action.Action
	PlayerId       string
	conn           net.Conn
}

func CreateFromConnection(conn net.Conn) User {
	var usr = User{make(chan action.Action), make(chan action.Action), uuid.New().String(), conn}
	go Send(usr)
	go Receive(usr)
	return usr
}

func Send(usr User) {
	encoder := json.NewEncoder(usr.conn)
	for {
		action := <-usr.SendChannel
		encoder.Encode(&action)
	}
}

func Receive(usr User) {
	decoder := json.NewDecoder(usr.conn)
	for {
		var action action.Action
		decoder.Decode(&action)
		usr.ReceiveChannel <- action
	}
}
