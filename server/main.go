package main

import (
	"fmt"

	"github.com/soyarielruiz/tdl-borbotones-go/server/gameManager"
)

func main() {
	var input string
	manager := gameManager.NewGameManager()
	go manager.Start()
	fmt.Print("Press q to quit server: \n")
	for input != "q" {
		fmt.Scanln(&input)
	}
	manager.Stop()
}
