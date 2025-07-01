package result

type ShuffleResult struct {
	ShuffleId    uint
	RoundResults []RoundResult
	Error        error
}

func NewShuffleResult(shuffleId uint, roundResults []RoundResult) ShuffleResult {
	return ShuffleResult{
		ShuffleId:    shuffleId,
		RoundResults: roundResults,
		Error:        nil,
	}
}

func NewShuffleResultWithError(shuffleId uint, err error) ShuffleResult {
	return ShuffleResult{
		ShuffleId: shuffleId,
		Error:     err,
	}
}

func (s ShuffleResult) GetNumRounds() int {
	return len(s.RoundResults)
}

func (s ShuffleResult) GetNumHands() int {
	numHands := 0
	for _, round := range s.RoundResults {
		numHands += round.GetNumHands()
	}
	return numHands
}

func (s ShuffleResult) GetBalance() int {
	balance := 0
	for _, round := range s.RoundResults {
		balance += round.GetBalance()
	}
	return balance
}
