package simulation

import (
	"errors"

	"github.com/jljl1337/blackjack-simulator/internal/blackjack"
	"github.com/jljl1337/blackjack-simulator/internal/core"
	"github.com/jljl1337/blackjack-simulator/internal/person"
	"github.com/jljl1337/blackjack-simulator/internal/result"
)

type ShuffleInput struct {
	ShuffleId uint
	Player    person.Player
	Dealer    person.Dealer
	Shoe      core.Shoe
	Rules     Rules
}

func PlayShuffleWorker(inputChan <-chan ShuffleInput, resultChan chan<- result.ShuffleResult) {
	for input := range inputChan {
		resultChan <- PlayShuffle(input)
	}
}

func PlayShuffle(input ShuffleInput) result.ShuffleResult {
	roundResults := make([]result.RoundResult, 0)

	shuffleId := input.ShuffleId
	player := input.Player
	dealer := input.Dealer
	shoe := input.Shoe
	rules := input.Rules

	for {
		if err := player.PlaceBet(); err != nil {
			return result.NewShuffleResultWithError(shuffleId, err)
		}

		dealer.DrawCard(shoe.Deal())
		dealer.DrawCard(shoe.Deal())
		player.DrawCard(shoe.Deal())
		player.DrawCard(shoe.Deal())

		playerHasBlackjack, err := player.CurrentHandIsBlackjack()
		if err != nil {
			return result.NewShuffleResultWithError(shuffleId, err)
		}

		dealerHasBlackjack := dealer.HasBlackjack()

		if dealerHasBlackjack {
			// Dealer has blackjack, check if player also has blackjack
			if !playerHasBlackjack {
				// Dealer wins
				player.LoseCurrentHand(1)
			}
			// Player also has blackjack, it's a push
			player.EndRound()
			dealer.EndRound()
			continue
		}

		// Player's turn
		for {
			actions, err := player.GetActions(dealer.GetUpCard())
			if err != nil {
				return result.NewShuffleResultWithError(shuffleId, err)
			}

			actionsAllowed, err := rules.GetActionsAllowed()
			if err != nil {
				return result.NewShuffleResultWithError(shuffleId, err)
			}

			selectedAction := blackjack.NA

			currentHandIsBlackjack, err := player.CurrentHandIsBlackjack()
			if err != nil {
				return result.NewShuffleResultWithError(shuffleId, err)
			}

			if currentHandIsBlackjack {
				// Skip selecting action if the player has blackjack
				selectedAction = blackjack.Blackjack
			} else {
				for _, action := range actions {
					if actionsAllowed[action] {
						selectedAction = action
						break
					}
				}
			}

			if selectedAction == blackjack.NA {
				return result.NewShuffleResultWithError(shuffleId, errors.New("no valid action selected"))
			}

			if selectedAction != blackjack.Blackjack {
				player.RecordAction(selectedAction)
			}

			playerLoseRatio := 0.0

			switch selectedAction {
			case blackjack.Blackjack:
				// If the dealer has blackjack, the player already pushed
			case blackjack.Hit:
				player.Hit(shoe.Deal())
				isBusted, err := player.CurrentHandIsBusted()
				if err != nil {
					return result.NewShuffleResultWithError(shuffleId, err)
				}
				if isBusted {
					// Player busts, dealer wins
					playerLoseRatio = 1.0 // Player loses the full bet
				} else {
					// Player hits, continue to next action
					continue
				}
			case blackjack.Stand:
			case blackjack.Double:
				if err := player.DoubleDown(shoe.Deal()); err != nil {
					return result.NewShuffleResultWithError(shuffleId, err)
				}
				isBusted, err := player.CurrentHandIsBusted()
				if err != nil {
					return result.NewShuffleResultWithError(shuffleId, err)
				}
				if isBusted {
					// Player busts, dealer wins
					playerLoseRatio = 1.0 // Player loses the full bet
				}
			case blackjack.Split:
				newCards := []core.Card{shoe.Deal(), shoe.Deal()}
				if err := player.Split(newCards); err != nil {
					return result.NewShuffleResultWithError(shuffleId, err)
				}
				// After splitting, player plays the new hand
				continue
			case blackjack.Surrender:
				playerLoseRatio = 0.5 // Player surrenders, loses half the bet
			}

			// Should only reach here if the current hand is ended
			if playerLoseRatio > 0 {
				// Player loses the hand, dealer wins
				player.LoseCurrentHand(playerLoseRatio)
			}

			if !player.HasNextHand() {
				break
			}

			// Player has more hands to play, continue to the next hand
			player.NextHand()
		}

		// Dealer's turn
		for dealer.NeedsToHit() {
			dealer.DrawCard(shoe.Deal())
		}

		dealerValue := dealer.GetHandValue()
		player.CalculateHandBet(dealerValue)

		roundResults = append(roundResults, result.NewRoundResult(
			dealer.GetHand(),
			player.GetHands(),
		))

		player.EndRound()
		dealer.EndRound()

		if shoe.NeedsShuffle() {
			// Finish this shuffle and start a new one
			return result.NewShuffleResult(shuffleId, roundResults)
		}
	}
}
