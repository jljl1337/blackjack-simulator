package person

import (
	"fmt"

	"github.com/jljl1337/blackjack-simulator/internal/core"
)

type Hand struct {
	cards []core.Card
}

func (h Hand) ValueString() string {
	value := h.Value()

	if h.IsSoft() {
		return fmt.Sprintf("S%d", value)
	}
	return fmt.Sprintf("H%d", value)
}

func (h Hand) PairString() (string, error) {
	if !h.IsPair() {
		return "", fmt.Errorf("hand is not a pair: %s", h.String())
	}

	return fmt.Sprintf("P%s", h.cards[0].ValueString()), nil
}

func (h Hand) Value() int {
	value := 0
	numAces := 0
	for _, card := range h.cards {
		if card.Rank == core.Ace {
			numAces++
		}
		lowValue, _ := card.Values()
		value += lowValue
	}

	for value <= 11 && numAces > 0 {
		value += 10
		numAces--
	}

	return value
}

func (h Hand) IsBlackjack() bool {
	return len(h.cards) == 2 && h.Value() == 21 && h.IsSoft()
}

func (h Hand) IsSoft() bool {
	lowVal, highVal := 0, 0
	for _, card := range h.cards {
		l, h := card.Values()
		lowVal += l
		highVal += h
	}
	return highVal != lowVal && highVal <= 21
}

func (h Hand) IsBusted() bool {
	return h.Value() > 21
}

func (h Hand) IsPair() bool {
	return len(h.cards) == 2 && h.cards[0].Rank == h.cards[1].Rank
}

func (h *Hand) AddCard(card core.Card) {
	h.cards = append(h.cards, card)
}

func (h Hand) String() string {
	str := ""
	for _, card := range h.cards {
		str += card.String() + ";"
	}
	return str[:len(str)-1]
}
