package simulation

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/jljl1337/blackjack-simulator/internal/blackjack"
	"github.com/jljl1337/blackjack-simulator/internal/core"
	"github.com/jljl1337/blackjack-simulator/internal/person"
)

type Simulator struct {
	seed        int64
	numShuffles int
}

func NewSimulator() *Simulator {
	seed := flag.Int64("seed", time.Now().UnixNano(), "Random seed for simulation")
	numShuffles := flag.Int("num-shuffles", 100000, "Number of shuffles to simulate")
	flag.Parse()

	fmt.Printf("Using seed: %d\n", *seed)
	return &Simulator{
		seed:        *seed,
		numShuffles: *numShuffles,
	}
}

func (s *Simulator) Run() {

	random := rand.New(rand.NewSource(s.seed))

	var balanceSum int64

	for i := range s.numShuffles {
		strategy, _ := blackjack.NewBasicStrategyS17()
		player := person.NewPlayer(strategy)
		dealer := person.NewDealer()
		rules := NewPlayRules()
		shoe := core.NewShoe(8, 0.75, random)
		result := PlayShuffle(i, *player, *dealer, *shoe, rules)
		if result.Error != nil {
			fmt.Printf("Error in shuffle %d: %v\n", i, result.Error)
			break
		}
		fmt.Printf("Shuffle %.7d: Played %.3d rounds with final balance of Balance: %d\n", i, result.NumRound, result.PlayerFinalBalance)

		balanceSum += int64(result.PlayerFinalBalance)
	}

	averageBalance := float64(balanceSum) / float64(s.numShuffles)
	fmt.Printf("Average balance after %d shuffles: %.2f\n", s.numShuffles, averageBalance)
	fmt.Printf("Final balance sum: %d\n", balanceSum)
}
