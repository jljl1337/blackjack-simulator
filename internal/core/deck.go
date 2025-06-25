package core

type Deck []Card

// NewDeck creates a standard 52-card deck
func NewDeck() Deck {
	deck := make(Deck, 52)
	i := 0
	for suit := Spades; suit <= Clubs; suit++ {
		for rank := Ace; rank <= King; rank++ {
			deck[i] = Card{Suit: suit, Rank: rank}
			i++
		}
	}
	return deck
}
