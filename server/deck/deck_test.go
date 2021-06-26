package deck

import (
	"github.com/soyarielruiz/tdl-borbotones-go/tools"
	"testing"
)

func TestDeck_GetCardFromDeck(t *testing.T) {
	deck := NewDeck()
	card := deck.GetCardFromDeck()
	if !isValidCard(card) {
		t.Errorf("Invalid card got")
	}
}

func isValidCard(card tools.Card) bool {
	return card.Number >= 0 && card.Number <= 9
}

