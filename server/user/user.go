package user

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/soyarielruiz/tdl-borbotones-go/tools"
)

type User struct {
	SendChannel    chan tools.Action
	ReceiveChannel chan tools.Action
	PlayerId       string
	conn           net.Conn
	Closed         bool
	CardsLeft      int
}

func (u User) String() string {
	return fmt.Sprintf("PlayerId:\"%s\", conn:\"%s\"", u.PlayerId, u.conn.RemoteAddr())
}

func (u *User) Close() {
	log.Printf("Closing user %s", u.PlayerId)
	u.Closed = true
	u.conn.Close()
	close(u.SendChannel)
}

func NewUser(conn net.Conn, nick string) *User {
	var usr = User{make(chan tools.Action), make(chan tools.Action), nick, conn, false, 3}
	go usr.Send()
	return &usr
}

func (usr *User) Send() {
	encoder := json.NewEncoder(usr.conn)
	for !usr.Closed {
		action,ok := <-usr.SendChannel
		if ok {
			log.Printf("Sending action to usr %s", action)
			encoder.Encode(&action)
		}

		if action.Command == tools.GAME_ENDED {
			usr.Close()
			log.Printf("Channel for user %s closed", usr.PlayerId)
		}
	}
}

func (usr *User) Receive() {
	decoder := json.NewDecoder(usr.conn)
	for !usr.Closed {
		var action tools.Action
		err := decoder.Decode(&action)
		if err != nil {
			if !usr.Closed {
				log.Printf("Sending close action from usr %s %b", usr.PlayerId, usr.Closed)
				usr.ReceiveChannel <- tools.Action{
					Command:  tools.EXIT,
					Card:     tools.Card{},
					PlayerId: usr.PlayerId,
					Message:  "",
					Cards:    nil,
				}
				return
			}
		} else {
			action.PlayerId = usr.PlayerId
			log.Printf("Receive action from usr %s.\n %s", usr, action)
			usr.ReceiveChannel <- action
		}
	}
}
