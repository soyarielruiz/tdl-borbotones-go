package user

import (
	"net"
)

type User struct{
	 Send_channel chan string
	 Recieve_channel chan string
}

func Send(sch chan string, connection net.Conn){
	 //TODO
}

func Recieve(rch chan string, connection net.Conn){
	 message:="empanada de ddl"
	 rch <- message
}