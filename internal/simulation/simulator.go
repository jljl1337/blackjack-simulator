package simulation

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
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
	numDecks    uint
	penetration float64
	csvFile     string
}

type Config struct {
	Seed        int64   `json:"seed"`
	NumShuffles uint    `json:"numShuffles"`
	NumDecks    uint    `json:"numDecks"`
	Penetration float64 `json:"penetration"`
}

func NewSimulator() *Simulator {
	configFile := flag.String("config", "config.json", "Path to configuration file")
	csvFile := flag.String("csv", "", "CSV file to export results to")

	flag.Parse()

	config, err := readConfig(*configFile)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", *configFile, err)
		return nil
	}

	fmt.Printf("Using seed: %d\n", config.Seed)
	return &Simulator{
		seed:        config.Seed,
		numShuffles: config.NumShuffles,
		numDecks:    config.NumDecks,
		penetration: config.Penetration,
		csvFile:     *csvFile,
	}
}

func readConfig(configFile string) (Config, error) {
	// Open the JSON file
	file, err := os.Open(configFile)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	// Read the file contents
	bytes, err := io.ReadAll(file)
	if err != nil {
		return Config{}, err
	}

	// Unmarshal JSON into the struct
	var config Config
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		switch e := err.(type) {
		case *json.UnmarshalTypeError:
			return Config{}, fmt.Errorf("field %s expected %s got %s", e.Field, e.Type, e.Value)
		default:
			return Config{}, err
		}
	}

	// Validate the config
	if config.NumShuffles <= 0 {
		return Config{}, fmt.Errorf("numShuffles must be set to a value greater than 0")
	}

	if config.NumDecks <= 0 {
		return Config{}, fmt.Errorf("numDecks must be set to a value greater than 0")
	}

	if config.Penetration <= 0 || config.Penetration > 1 {
		return Config{}, fmt.Errorf("penetration must be set to a value larger than 0 and at most 1")
	}

	// Set default values if not provided
	if config.Seed == 0 {
		config.Seed = time.Now().UnixNano()
	}

	return config, nil
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
		shoe := core.NewShoe(s.numDecks, s.penetration, random)
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
