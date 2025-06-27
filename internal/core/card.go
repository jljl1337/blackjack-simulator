package core

import "fmt"

type Card struct {
	Suit Suit
	Rank Rank
}

// ValueString returns the string representation of the card's value.
// For Ace, it returns "A".
func (c Card) ValueString() string {
	if c.Rank == Ace {
		return "A"
	}

	value, _ := c.Values()

	return fmt.Sprintf("%d", value)
}

// Values returns the possible values of the card.
func (c Card) Values() (int, int) {
	switch c.Rank {
	case Ace:
		return 1, 11
	case King, Queen, Jack, Ten:
		return 10, 10
	default:
		return int(c.Rank), int(c.Rank)
	}
}

func (c Card) String() string {
	ranks := map[Rank]string{
		Ace:   "A",
		Two:   "2",
		Three: "3",
		Four:  "4",
		Five:  "5",
		Six:   "6",
		Seven: "7",
		Eight: "8",
		Nine:  "9",
		Ten:   "10",
		Jack:  "J",
		Queen: "Q",
		King:  "K",
	}
	suits := map[Suit]string{
		Spades:   "♠",
		Hearts:   "♥",
		Diamonds: "♦",
		Clubs:    "♣",
	}
	return fmt.Sprintf("%s%s", ranks[c.Rank], suits[c.Suit])
}
