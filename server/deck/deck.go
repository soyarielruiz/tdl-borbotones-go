package deck

import (
	"math/rand"
	"time"

	"github.com/soyarielruiz/tdl-borbotones-go/server/stack"
	"github.com/soyarielruiz/tdl-borbotones-go/tools"
)

type Deck struct {
	deckStack    *stack.Stack
	discardStack *stack.Stack
}

func NewDeck() *Deck {
	deck := Deck{stack.New(), stack.New()}
	numbers := [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	suits := [4]tools.Suit{tools.RED, tools.GREEN, tools.BLUE, tools.YELLOW}
	for i := 1; i < 3; i++ {
		for _, s := range suits {
			for _, n := range numbers {
				deck.discardStack.Push(tools.Card{n, s})
			}
		}
	}
	deck.shuffle()
	return &deck
}

func (deck *Deck) shuffle() {
	lastCard, _ := deck.discardStack.Pop()
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(deck.discardStack.Size(), func(i, j int) { deck.discardStack.Swap(i, j) })
	deck.deckStack.PushAll(*deck.discardStack)
	deck.discardStack.Clear()
	deck.discardStack.Push(lastCard)
}

func (deck *Deck) GetCardFromDeck() tools.Card {
	c, empty := deck.deckStack.Pop()
	if empty {
		deck.shuffle()
		return deck.GetCardFromDeck()
	}
	return c.(tools.Card)
}

func (deck *Deck) GetCardsFromDeck(n int) []tools.Card {
	result := make([]tools.Card, n)
	for i := 0; i < n; i++ {
		result[i] = deck.GetCardFromDeck()
	}
	return result
}

func (deck *Deck) PutCard(card tools.Card) {
	deck.discardStack.Push(card)
}

func (deck *Deck) GetFrontCard() tools.Card {
	result, _ := deck.discardStack.Front()
	return result.(tools.Card)
}
