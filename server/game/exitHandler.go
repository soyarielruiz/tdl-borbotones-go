package game

import (
	"github.com/soyarielruiz/tdl-borbotones-go/tools"
	"log"
)

type ExitHandler struct{}

func (t ExitHandler) Handle(action tools.Action, game *Game) {
	log.Printf("Handle exit from usr %s", action.PlayerId)
	game.Users[action.PlayerId].Close()
	delete(game.Users, action.PlayerId)
	action.Message= "Player " + action.PlayerId + " has disconnected"
	game.SendToAll(&action)
	previousUser := game.Tur.CurrentUser()
	game.Tur.Remove(action.PlayerId)
	if previousUser == action.PlayerId {
		game.Users[game.Tur.CurrentUser()].SendChannel <- tools.Action{
			Command:  tools.TURN_ASSIGNED,
			Card:     game.Deck.GetFrontCard(),
			PlayerId: game.Tur.CurrentUser(),
			Message:  "It's your turn to play!",
			Cards:    nil,
		}
	}
	if len(game.Users) < 2 {
		log.Printf("No enough players to continue the game %d.", game.GameNumber)
		game.Ended = true
		game.SendToAll(&tools.Action{
			Command:  tools.GAME_ENDED,
			Card:     tools.Card{},
			PlayerId: "",
			Message:  "Game ended! No enough players",
			Cards:    nil,
		})

	}
}
