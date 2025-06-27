package blackjack

import (
	"github.com/jljl1337/blackjack-simulator/internal/core"
)

type Strategy interface {
	GetActions(playerHand core.Hand, dealerUpCard core.Card) ([]Action, error)
}
