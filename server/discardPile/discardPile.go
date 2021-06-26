package discardPile

import "github.com/soyarielruiz/tdl-borbotones-go/tools"

type DiscardPile struct {
	LastCard tools.Card
}

func NewDiscardPile(firstCard tools.Card) *DiscardPile {
	return &DiscardPile{firstCard}
}

func (d *DiscardPile) PutCard(card tools.Card) {
	d.LastCard = card
}

func (d *DiscardPile) GetLastCard() tools.Card {
	return d.LastCard
}
