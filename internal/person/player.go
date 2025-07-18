package person

import (
	"errors"

	"github.com/jljl1337/blackjack-simulator/internal/blackjack"
	"github.com/jljl1337/blackjack-simulator/internal/core"
)

type Player struct {
	currentHand int
	hands       []*PlayerHand
	strategy    blackjack.Strategy
}

func NewPlayer(strategy blackjack.Strategy) *Player {
	return &Player{
		hands:    []*PlayerHand{NewPlayerHand()},
		strategy: strategy,
	}
}

func (p *Player) PlaceBet() error {
	const betAmount = 100 // TODO: card counting
	currentHand, err := p.getCurrentHand()
	if err != nil {
		return err
	}

	currentHand.PlaceBet(betAmount)

	return nil
}

// CalculateHandBet calculates the final value of each bet at the end of the
// round, based on the dealer's hand value and the player's hand value.
func (p *Player) CalculateHandBet(dealerValue int) {
	for i, hand := range p.hands {
		if hand.GetBet() != hand.GetBetPlaced() {
			// This hand is either busted or surrendered
		} else if hand.IsBlackjack() {
			// Note that if the dealer has a blackjack, the bet is calculated at
			// the beginning of the round, so we don't need to check for that.
			p.winHand(i, 1.5)
		} else if dealerValue > 21 || hand.Value() > dealerValue {
			p.winHand(i, 1.0)
		} else if hand.Value() < dealerValue {
			p.loseHand(i, 1.0)
		}

		// If the values are equal, the bet remains the same (push)
	}
}

func (p *Player) DrawCard(card core.Card) error {
	currentHand, err := p.getCurrentHand()
	if err != nil {
		return err
	}

	if currentHand.IsBusted() {
		return errors.New("cannot draw a card for a busted hand")
	}

	currentHand.AddCard(card)
	return nil
}

func (p Player) CurrentHandIsBlackjack() (bool, error) {
	currentHand, err := p.getCurrentHand()
	if err != nil {
		return false, err
	}

	return currentHand.IsBlackjack(), nil
}

func (p Player) CurrentHandIsBusted() (bool, error) {
	currentHand, err := p.getCurrentHand()
	if err != nil {
		return false, err
	}

	return currentHand.IsBusted(), nil
}

func (p Player) GetCurrentHandSize() (int, error) {
	currentHand, err := p.getCurrentHand()
	if err != nil {
		return 0, err
	}

	return len(currentHand.cards), nil
}

func (p *Player) GetActions(dealerUpCard core.Card) ([]blackjack.Action, error) {
	currentHand, err := p.getCurrentHand()
	if err != nil {
		return nil, err
	}

	// Use the strategy to get the actions for the current hand
	actions, err := p.strategy.GetActions(currentHand, dealerUpCard)
	if err != nil {
		return nil, err
	}

	return actions, nil
}

func (p *Player) RecordAction(action blackjack.Action) error {
	currentHand, err := p.getCurrentHand()
	if err != nil {
		return err
	}

	currentHand.AddAction(action)
	return nil
}

func (p *Player) Hit(newCard core.Card) error {
	currentHand, err := p.getCurrentHand()
	if err != nil {
		return err
	}

	if currentHand.IsBusted() {
		return errors.New("cannot hit a busted hand")
	}

	currentHand.AddCard(newCard)
	return nil
}

func (p *Player) DoubleDown(newCard core.Card) error {
	currentHand, err := p.getCurrentHand()
	if err != nil {
		return err
	}
	if currentHand.IsBusted() {
		return errors.New("cannot double down on a busted hand")
	}

	// Double the bet
	currentHand.PlaceBet(currentHand.GetBet() * 2)

	// Add the new card to the current hand
	currentHand.AddCard(newCard)

	return nil
}

func (p *Player) Split(newCards []core.Card) error {
	if len(newCards) != 2 {
		return errors.New("split requires exactly two new card")
	}

	currentHand, err := p.getCurrentHand()
	if err != nil {
		return err
	}

	if len(currentHand.cards) != 2 {
		return errors.New("cannot split a hand that does not have exactly two cards")
	}

	// Create a new hand
	newHand := NewPlayerHand()
	p.hands = append(p.hands, newHand)

	// Adjust the bet for the new hand
	newHand.PlaceBet(currentHand.GetBet())

	// Move the second card from the current hand to the new hand
	newHand.AddCard(currentHand.cards[1])
	currentHand.cards = currentHand.cards[:1]

	// Add the new cards to the current hand
	currentHand.AddCard(newCards[0])
	newHand.AddCard(newCards[1])

	return nil
}

func (p *Player) WinCurrentHand(winRatio float64) error {
	return p.winHand(p.currentHand, winRatio)
}

func (p *Player) LoseCurrentHand(loseRatio float64) error {
	return p.loseHand(p.currentHand, loseRatio)
}

func (p *Player) winHand(index int, winRatio float64) error {
	return p.adjustHandBet(index, 1.0+winRatio)
}

func (p *Player) loseHand(index int, loseRatio float64) error {
	return p.adjustHandBet(index, 1.0-loseRatio)
}

func (p *Player) adjustHandBet(index int, ratio float64) error {
	hand, err := p.getHand(index)
	if err != nil {
		return err
	}

	hand.AdjustBetByRatio(ratio)

	return nil
}

func (p *Player) getCurrentHand() (*PlayerHand, error) {
	return p.getHand(p.currentHand)
}

func (p *Player) getHand(index int) (*PlayerHand, error) {
	if len(p.hands) == 0 {
		return nil, errors.New("no hands available")
	}

	if index < 0 || index >= len(p.hands) {
		return nil, errors.New("invalid hand index")
	}

	return p.hands[index], nil
}

func (p *Player) NextHand() error {
	if !p.HasNextHand() {
		return errors.New("no next hand available")
	}

	p.currentHand++
	return nil
}

func (p *Player) HasNextHand() bool {
	return p.currentHand < len(p.hands)-1
}

func (p Player) SplitAce() bool {
	return len(p.hands) > 1 && p.hands[0].cards[0].Rank == core.Ace
}

func (p Player) GetNumHands() int {
	return len(p.hands)
}

func (p Player) GetHands() []*PlayerHand {
	return p.hands
}

// EndRound resets the player's state for a new round
func (p *Player) EndRound() {
	p.currentHand = 0
	p.hands = []*PlayerHand{NewPlayerHand()}
}
