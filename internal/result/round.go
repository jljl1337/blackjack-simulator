package result

import (
	"github.com/jljl1337/blackjack-simulator/internal/person"
)

type RoundResult struct {
	DealerHand  person.Hand
	PlayerHands []person.PlayerHand
}

func NewRoundResult(dealerHand person.Hand, playerHands []*person.PlayerHand) RoundResult {
	hands := make([]person.PlayerHand, len(playerHands))
	for i, hand := range playerHands {
		hands[i] = *hand
	}
	return RoundResult{
		DealerHand:  dealerHand,
		PlayerHands: hands,
	}
}

func (r RoundResult) GetNumHands() int {
	return len(r.PlayerHands)
}
