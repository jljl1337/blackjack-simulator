package person

import "github.com/jljl1337/blackjack-simulator/internal/core"

type Dealer struct {
	hand Hand
}

func NewDealer() *Dealer {
	return &Dealer{
		hand: Hand{},
	}
}

func (d *Dealer) DrawCard(card core.Card) {
	d.hand.cards = append(d.hand.cards, card)
}

func (d Dealer) HasBlackjack() bool {
	return d.hand.IsBlackjack()
}

func (d Dealer) NeedsToHit() bool {
	// Dealer hits on 16 or less, stands on 17 or more
	return d.hand.Value() < 17
}

func (d Dealer) GetUpCard() core.Card {
	return d.hand.cards[0]
}

func (d Dealer) GetHandValue() int {
	return d.hand.Value()
}

func (d *Dealer) EndRound() {
	d.hand = Hand{}
}
