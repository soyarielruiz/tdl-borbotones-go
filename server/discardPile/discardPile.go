package discardPile

import "github.com/soyarielruiz/tdl-borbotones-go/tools"

type DiscardPile struct {
	CurrentCard tools.Card
	LastCard    tools.Card
}

func (d DiscardPile) Init() {

}
func (d *DiscardPile) PutCard(card tools.Card) {
	//TODO
}

func (d *DiscardPile) GetLastCard() tools.Card {
	//TODO
	return tools.Card{}
}
