package main

import (
	"errors"
	"fmt"
	"github.com/awesome-gocui/gocui"
	"github.com/soyarielruiz/tdl-borbotones-go/client/translator"
	"github.com/soyarielruiz/tdl-borbotones-go/client/game"
	"github.com/soyarielruiz/tdl-borbotones-go/tools"
	"log"
	"os"
	"strconv"
	"github.com/soyarielruiz/tdl-borbotones-go/client/view"
)

type LobbyOption struct{
	Option []int `json:"option"`
}

func main() {

	game:=lobby()

	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen

	g.SetManagerFunc(view.Layout)

	// Bind enter key to input to send new messages.
	err = g.SetKeybinding("jugador", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, iv *gocui.View) error {
		iv.Autoscroll = true
		iv.Editable = true
		// Read buffer from the beginning.
		iv.Rewind()

		// Send message if text was entered.
		if len(iv.Buffer()) >= 2 {

			messageToUse := string(iv.Buffer())
			messageToSend, err := translator.CreateAnAction(messageToUse, g)
			if err != nil {
				out, _ := g.View("mano")
				fmt.Fprintf(out, "Error al crear la accion, probar nuevamente \n")
			}

			if translator.MustLeave(messageToSend) {
				return view.Quit(g, iv)
			}

			if translator.GameWasEnded(messageToSend) {
				out, _ := g.View("mano")
				fmt.Fprintf(out, "Juego Finalizado \n")
				return nil
			}

			if translator.HaveActionToSend(messageToSend) {
				err = game.Encoder.Encode(&messageToSend)
				checkError(err)
			}

			//get cursor position
			_, y := iv.Cursor()

			//adding a visual enter
			w := y + 1

			err = iv.SetCursorUnrestricted(0, w)
			if err != nil {
				log.Println("Failed to set cursor:", err)
			}

			iv.Clear()
			return err
		}
		return nil
	})

	if err != nil {
		log.Println("Cannot bind the enter key:", err)
	}

	if err := view.InitKeybindings(g); err != nil {
		log.Fatalln(err)
	}

	// receiving from server
	go receivingData(g,game)

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
	game:=game.NewGame()
	fmt.Println("* * * * * * * * * *")
	fmt.Println("* WELCOME TO GUNO *")
	fmt.Println("* * * * * * * * * *")
	for {
		input:=initialOption()
		option :=LobbyOption{[]int{input}}
		game.Encoder.Encode(&option)
		if input==2{
			var games LobbyOption
			game.Decoder.Decode(&games)
			if (len(games.Option)==0){
				fmt.Println("There are no current games available.")
				continue
			}
			input=gameOption(games)
			option2 :=LobbyOption{[]int{input}}
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
		if option!=1 && option!=2 || err!=nil {
			fmt.Println("Wrong option. Please type 1 or 2: ")
		}else{
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
		if err!=nil || !checkGameId(games,option) {
			fmt.Println("Wrong option")
		}else{
			break
		}
	}
	return option
}

func checkGameId(games LobbyOption,option int) bool{
	for _, gameId := range games.Option {
		if gameId == option {
			return true
		}
	}
	return false
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
