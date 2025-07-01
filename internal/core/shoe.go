package core

import (
	"math/rand"
)

// Shoe represents multiple decks of cards used in Blackjack
type Shoe struct {
	cards       []Card
	penetration float64
	numDecks    uint
}

// NewShoe creates a shoe with a specified number of decks
func NewShoe(numDecks uint, penetration float64, rand *rand.Rand) *Shoe {
	s := &Shoe{numDecks: numDecks, penetration: penetration}
	for range numDecks {
		s.cards = append(s.cards, NewDeck()...)
	}

	rand.Shuffle(len(s.cards), func(i, j int) {
		s.cards[i], s.cards[j] = s.cards[j], s.cards[i]
	})

	return s
}

// Deal deals a card from the shoe
func (s *Shoe) Deal() Card {
	if len(s.cards) == 0 {
		panic("Shoe is empty!")
	}
	card := s.cards[0]
	s.cards = s.cards[1:]
	return card
}

// NeedsShuffle checks if the shoe needs to be shuffled based on penetration
func (s *Shoe) NeedsShuffle() bool {
	return float64(len(s.cards)) < float64(s.numDecks*52)*(1.0-s.penetration)
}
