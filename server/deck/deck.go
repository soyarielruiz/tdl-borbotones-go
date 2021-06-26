package deck

import (
	"github.com/soyarielruiz/tdl-borbotones-go/server/stack"
	"github.com/soyarielruiz/tdl-borbotones-go/tools"
	"math/rand"
	"time"
)

type Deck struct {
	stack *stack.Stack
}

func NewDeck() *Deck {
	deck := Deck{stack.New()}
	numbers := [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	suits := [4]tools.Suit{tools.RED, tools.GREEN, tools.BLUE, tools.YELLOW}
	for _, s := range suits {
		for _, n := range numbers {
			deck.stack.Push(tools.Card{n, s})
		}
	}
	deck.shuffle()
	return &deck;
}

func (deck *Deck) shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(deck.stack.Size(), func(i, j int) { deck.stack.Swap(i, j) })

}

func (deck *Deck) GetCard() tools.Card {
	//TODO
	return tools.Card{}
}

func (deck *Deck) GetCards(n int) []tools.Card {
	//TODO
	return []tools.Card{}
}
