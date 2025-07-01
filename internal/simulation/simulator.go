package simulation

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"runtime"
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
	numWorkers  uint
	strategy    blackjack.Strategy
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
	numWorkers := flag.Uint("num-workers", 0, "Number of workers to use for concurrent processing")

	flag.Parse()

	if *numWorkers == 0 {
		*numWorkers = uint(runtime.NumCPU())
	}

	config, err := readConfig(*configFile)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", *configFile, err)
		return nil
	}

	fmt.Printf("Using seed: %d\n", config.Seed)
	fmt.Printf("Number of workers: %d\n", *numWorkers)

	strategy, err := blackjack.NewBasicStrategyS17()
	if err != nil {
		fmt.Printf("Error creating strategy: %v\n", err)
		return nil
	}

	return &Simulator{
		seed:        config.Seed,
		numShuffles: config.NumShuffles,
		numDecks:    config.NumDecks,
		penetration: config.Penetration,
		csvFile:     *csvFile,
		numWorkers:  *numWorkers,
		strategy:    strategy,
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

	inputChan := make(chan ShuffleInput, s.numWorkers)
	resultChan := make(chan result.ShuffleResult, s.numWorkers)

	// Start consumer workers
	for range s.numWorkers {
		go PlayShuffleWorker(inputChan, resultChan)
	}

	shuffleId := uint(0)

	// Send numWorkers shuffle inputs to the input channel
	for range s.numWorkers {
		s.sendInput(inputChan, shuffleId, random)
		shuffleId++
	}

	// Collect results from the result channel
	shuffleResults := make([]result.ShuffleResult, s.numWorkers)

	count := uint(0)
	countedShuffles := uint(0)

out:
	for {
		shuffleResult := <-resultChan
		if shuffleResult.Error != nil {
			fmt.Printf("Error in shuffle %d: %v\n", shuffleResult.ShuffleId, shuffleResult.Error)
			break
		}

		// fmt.Printf("Shuffle %d: Played %d rounds with final balance of: %d\n", shuffleResult.ShuffleId, shuffleResult.NumRounds, shuffleResult.Balance)

		shuffleResults[shuffleResult.ShuffleId] = shuffleResult

		if shuffleResult.ShuffleId == countedShuffles {
			for i := countedShuffles; i < uint(len(shuffleResults)); i++ {
				if shuffleResults[i].NumRounds <= 0 {
					break
				}

				countedShuffles++
				count++

				if count >= s.numShuffles {
					fmt.Printf("finished %d shuffles\n", countedShuffles)
					close(inputChan)
					break out
				}
			}
		}

		shuffleResults = append(shuffleResults, result.NewShuffleResult(0, nil))
		s.sendInput(inputChan, shuffleId, random)
		shuffleId++
	}

	shuffleResults = shuffleResults[:countedShuffles]

	var balanceSum int64
	for _, result := range shuffleResults {
		balanceSum += int64(result.Balance)
	}
	averageBalance := float64(balanceSum) / float64(countedShuffles)

	fmt.Printf("Average balance after %d shuffles: %.2f\n", countedShuffles, averageBalance)
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

func (s *Simulator) sendInput(inputChan chan<- ShuffleInput, shuffleId uint, random *rand.Rand) {
	player := person.NewPlayer(s.strategy)
	dealer := person.NewDealer()
	rules := NewPlayRules()
	shoe := core.NewShoe(s.numDecks, s.penetration, random)

	input := ShuffleInput{
		ShuffleId: shuffleId,
		Player:    *player,
		Dealer:    *dealer,
		Shoe:      *shoe,
		Rules:     rules,
	}

	inputChan <- input
}
