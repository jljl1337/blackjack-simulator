package person

import (
	"github.com/jljl1337/blackjack-simulator/internal/blackjack"
	"github.com/jljl1337/blackjack-simulator/internal/core"
)

type PlayerHand struct {
	Hand
	betPlaced int
	bet       int
	actions   []blackjack.Action
}

func NewPlayerHand() *PlayerHand {
	return &PlayerHand{
		Hand: Hand{
			cards: []core.Card{},
		},
		bet: 0,
	}
}

func (ph *PlayerHand) PlaceBet(amount int) {
	ph.betPlaced = amount
	ph.bet = amount
}

func (ph *PlayerHand) AdjustBetByRatio(ratio float64) {
	ph.bet = int(float64(ph.bet) * ratio)
}

func (ph PlayerHand) GetBetPlaced() int {
	return ph.betPlaced
}

func (ph PlayerHand) GetBet() int {
	return ph.bet
}

func (ph *PlayerHand) AddAction(action blackjack.Action) {
	ph.actions = append(ph.actions, action)
}

func (ph PlayerHand) GetActions() []blackjack.Action {
	return ph.actions
}
