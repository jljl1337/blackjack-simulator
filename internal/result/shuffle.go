package result

type ShuffleResult struct {
	ShuffleId    int
	RoundResults []RoundResult
	Error        error
}

func NewShuffleResult(shuffleId int, roundResults []RoundResult) ShuffleResult {
	return ShuffleResult{
		ShuffleId:    shuffleId,
		RoundResults: roundResults,
		Error:        nil,
	}
}

func NewShuffleResultWithError(shuffleId int, err error) ShuffleResult {
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
