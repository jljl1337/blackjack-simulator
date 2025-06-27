package core

// Hand represents a player's or dealer's hand in a game of blackjack.
type Hand interface {
	// Value returns the total value of the hand.
	ValueString() string
	// PairString returns a string representation of the hand if it is a pair.
	// If the hand is not a pair, it returns an error.
	PairString() (string, error)
	// Value returns the total value of the hand.
	Value() int
	// IsBlackjack checks if the hand is a blackjack (21 with two cards).
	IsBlackjack() bool
	// IsSoft checks if the hand is a soft hand (contains an Ace counted as
	// 11).
	IsSoft() bool
	// IsBusted checks if the hand is busted (value exceeds 21).
	IsBusted() bool
	// IsPair checks if the hand is a pair (two cards of the same rank).
	IsPair() bool
}
