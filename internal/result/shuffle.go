package result

type ShuffleResult struct {
	ShuffleId    uint
	RoundResults []RoundResult
	NumRounds    int
	NumHands     int
	Balance      int
	Error        error
}

func NewShuffleResult(shuffleId uint, roundResults []RoundResult) ShuffleResult {
	numHands := 0
	balance := 0
	for _, round := range roundResults {
		numHands += round.NumHands
		balance += round.Balance
	}

	return ShuffleResult{
		ShuffleId:    shuffleId,
		RoundResults: roundResults,
		NumRounds:    len(roundResults),
		NumHands:     numHands,
		Balance:      balance,
		Error:        nil,
	}
}

func NewShuffleResultWithError(shuffleId uint, err error) ShuffleResult {
	return ShuffleResult{
		ShuffleId: shuffleId,
		Error:     err,
	}
}
