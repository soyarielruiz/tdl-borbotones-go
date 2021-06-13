package game

import (
	"fmt"

	"github.com/soyarielruiz/tdl-borbotones-go/server/user"
)


func Start(ch chan user.User, game_number int) {
	 fmt.Printf("new game number: %d\n",game_number)
	 for {
		user := <-ch
		fmt.Printf("new player in game: %d ",game_number)
		fmt.Println("says: ", <-user.Recieve_channel)
	 } 
}
