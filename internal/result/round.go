package result

import (
	"github.com/jljl1337/blackjack-simulator/internal/person"
)

type RoundResult struct {
	DealerHand  person.Hand
	PlayerHands []person.PlayerHand
	NumHands    int
	Balance     int
}

func NewRoundResult(dealerHand person.Hand, playerHands []*person.PlayerHand) RoundResult {
	numHands := len(playerHands)
	hands := make([]person.PlayerHand, numHands)
	for i, hand := range playerHands {
		hands[i] = *hand
	}

	balance := 0
	for _, hand := range playerHands {
		balance += hand.GetBet() - hand.GetBetPlaced()
	}

	return RoundResult{
		DealerHand:  dealerHand,
		PlayerHands: hands,
		NumHands:    numHands,
		Balance:     balance,
	}
}
