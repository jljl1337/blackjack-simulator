package simulation

import (
	"github.com/jljl1337/blackjack-simulator/internal/blackjack"
	"github.com/jljl1337/blackjack-simulator/internal/core"
)

type ShuffleResult struct {
	ShuffleId          int
	NumRound           int
	PlayerFinalBalance int
	Error              error
}

func NewShuffleResult(shuffleId int, numRound int, playerFinalBalance int) ShuffleResult {
	return ShuffleResult{
		ShuffleId:          shuffleId,
		NumRound:           numRound,
		PlayerFinalBalance: playerFinalBalance,
		Error:              nil,
	}
}

func NewShuffleResultWithError(shuffleId int, err error) ShuffleResult {
	return ShuffleResult{
		ShuffleId:          shuffleId,
		NumRound:           0,
		PlayerFinalBalance: 0,
		Error:              err,
	}
}

type HandResult struct {
	Hand                 core.Hand
	BetBeforeCalculation int
	BetAfterCalculation  int
	Actions              []blackjack.Action
}
