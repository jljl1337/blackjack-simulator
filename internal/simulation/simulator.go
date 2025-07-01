package simulation

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/jljl1337/blackjack-simulator/internal/blackjack"
	"github.com/jljl1337/blackjack-simulator/internal/core"
	"github.com/jljl1337/blackjack-simulator/internal/exporter"
	"github.com/jljl1337/blackjack-simulator/internal/person"
	"github.com/jljl1337/blackjack-simulator/internal/result"
)

type Simulator struct {
	seed        int64
	numShuffles uint
	csvFile     string
}

func NewSimulator() *Simulator {
	seed := flag.Int64("seed", time.Now().UnixNano(), "Random seed for simulation")
	numShuffles := flag.Uint("num-shuffles", 100000, "Number of shuffles to simulate")
	csvFile := flag.String("csv", "", "CSV file to export results to")

	flag.Parse()

	fmt.Printf("Using seed: %d\n", *seed)
	return &Simulator{
		seed:        *seed,
		numShuffles: *numShuffles,
		csvFile:     *csvFile,
	}
}

func (s *Simulator) Run() {
	random := rand.New(rand.NewSource(s.seed))

	var balanceSum int64
	var shuffleResults []result.ShuffleResult

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

		roundBalance := result.GetBalance()

		fmt.Printf("Shuffle %d: Played %d rounds with final balance of: %d\n", i, result.GetNumRounds(), roundBalance)

		balanceSum += int64(roundBalance)

		shuffleResults = append(shuffleResults, result)
	}

	averageBalance := float64(balanceSum) / float64(s.numShuffles)
	fmt.Printf("Average balance after %d shuffles: %.2f\n", s.numShuffles, averageBalance)
	fmt.Printf("Final balance sum: %d\n", balanceSum)

	if s.csvFile != "" {
		csvExporter := exporter.NewCSVExporter(s.csvFile)
		fmt.Printf("Exporting results to CSV...\n")
		if err := csvExporter.Export(shuffleResults); err != nil {
			fmt.Printf("Error exporting results to CSV: %v\n", err)
			return
		}
	}
}
