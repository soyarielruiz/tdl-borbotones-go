package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/awesome-gocui/gocui"
	"github.com/soyarielruiz/tdl-borbotones-go/client/game"
	"github.com/soyarielruiz/tdl-borbotones-go/client/view"
	"github.com/soyarielruiz/tdl-borbotones-go/tools"
)

type LobbyOption struct {
	Option []int `json:"option"`
}

func main() {

	game := lobby()

	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(view.Layout)

	if err := view.InitKeybindings(g, game); err != nil {
		log.Fatalln(err)
	}

	// receiving from server
	go receivingData(g, game)

	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Panicln(err)
	}
}

func receivingData(g *gocui.Gui, game *game.Game) {
	for {
		if err := view.ReceiveMsgFromGame(g, game); err != nil {
			log.Fatalln(err)
		}
	}
}

func lobby() *game.Game {
	game := game.NewGame()
	fmt.Println("* * * * * * * * * *")
	fmt.Println("* WELCOME TO GUNO *")
	fmt.Println("* * * * * * * * * *")
	for {
		input := initialOption()
		option := LobbyOption{[]int{input}}
		game.Encoder.Encode(&option)
		if input == 2 {
			var games LobbyOption
			game.Decoder.Decode(&games)
			if len(games.Option) == 0 {
				fmt.Println("There are no current games available.")
				continue
			}
			input = gameOption(games)
			option2 := LobbyOption{[]int{input}}
			game.Encoder.Encode(&option2)
		}
		break
	}
	fmt.Println("Waiting for new members to start")
	var start tools.Action
	game.Decoder.Decode(&start)
	return game
}

func initialOption() int {
	var option int
	var s string
	for {
		fmt.Println("1:Start new game\n2:Join game")
		_, err := fmt.Scan(&s)
		option, err = strconv.Atoi(s)
		if option != 1 && option != 2 || err != nil {
			fmt.Println("Wrong option. Please type 1 or 2: ")
		} else {
			break
		}
	}
	return option
}

func gameOption(games LobbyOption) int {
	var option int
	var s string
	for {
		fmt.Println("Choose game number:")
		fmt.Println(games.Option)
		_, err := fmt.Scan(&s)
		option, err = strconv.Atoi(s)
		if err != nil || !checkGameId(games, option) {
			fmt.Println("Wrong option")
		} else {
			break
		}
	}
	return option
}

func checkGameId(games LobbyOption, option int) bool {
	for _, gameId := range games.Option {
		if gameId == option {
			return true
		}
	}
	return false
}
