package user

type User struct {
	Id        uint
	Send_chan chan string
	Recv_chan chan string
}
