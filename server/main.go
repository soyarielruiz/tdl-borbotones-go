package main

import (
	"fmt"
	"os"

	"github.com/soyarielruiz/tdl-borbotones-go/server/gameManager"
)

func main() {
	 var input string
	 go gameManager.Start()
	 fmt.Print("Press q to quit server: \n")
	 for input != "q" {
		_, err := fmt.Scanln(&input)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
	 }
}
