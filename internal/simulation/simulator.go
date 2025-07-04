package simulation

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
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
	seed                int64
	numShuffles         uint
	numRounds           uint
	numHands            uint
	numDecks            uint
	penetration         float64
	csvFile             string
	numWorkers          uint
	verbose             bool
	strategy            blackjack.Strategy
	doubleAfterSplit    bool
	hitAfterSplitAce    bool
	splitAfterSplitAce  bool
	doubleAfterSplitAce bool
	maxNumHands         int
}

type Config struct {
	Seed                int64   `json:"seed"`
	NumShuffles         uint    `json:"numShuffles"`
	NumRounds           uint    `json:"numRounds"`
	NumHands            uint    `json:"numHands"`
	NumDecks            uint    `json:"numDecks"`
	Penetration         float64 `json:"penetration"`
	DoubleAfterSplit    bool    `json:"doubleAfterSplit"`
	HitAfterSplitAce    bool    `json:"hitAfterSplitAce"`
	SplitAfterSplitAce  bool    `json:"splitAfterSplitAce"`
	DoubleAfterSplitAce bool    `json:"doubleAfterSplitAce"`
	MaxNumHands         *int    `json:"maxNumHands"`
}

func NewSimulator() (*Simulator, error) {
	configFile := flag.String("config", "config.json", "Path to configuration file")
	csvFile := flag.String("csv", "", "CSV file to export results to")
	numWorkers := flag.Uint("num-workers", 0, "Number of workers to use for concurrent processing")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")

	flag.Parse()

	if *verbose {
		log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	}

	if *numWorkers == 0 {
		*numWorkers = uint(runtime.NumCPU())
	}

	config, err := readConfig(*configFile)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %w", *configFile, err)
	}

	log.Printf("Using seed: %d\n", config.Seed)
	log.Printf("Number of workers: %d\n", *numWorkers)

	strategy, err := blackjack.NewBasicStrategyS17()
	if err != nil {
		return nil, fmt.Errorf("error creating strategy: %w", err)
	}

	maxNumHands := 4
	if config.MaxNumHands != nil {
		maxNumHands = *config.MaxNumHands
	}

	return &Simulator{
		seed:                config.Seed,
		numShuffles:         config.NumShuffles,
		numDecks:            config.NumDecks,
		numRounds:           config.NumRounds,
		numHands:            config.NumHands,
		penetration:         config.Penetration,
		csvFile:             *csvFile,
		numWorkers:          *numWorkers,
		verbose:             *verbose,
		strategy:            strategy,
		doubleAfterSplit:    config.DoubleAfterSplit,
		hitAfterSplitAce:    config.HitAfterSplitAce,
		splitAfterSplitAce:  config.SplitAfterSplitAce,
		doubleAfterSplitAce: config.DoubleAfterSplitAce,
		maxNumHands:         maxNumHands,
	}, nil
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
	conditionCount := 0
	if config.NumShuffles > 0 {
		conditionCount++
	}
	if config.NumRounds > 0 {
		conditionCount++
	}
	if config.NumHands > 0 {
		conditionCount++
	}

	if conditionCount != 1 {
		return Config{}, fmt.Errorf("exactly one of numShuffles, numRounds, or numHands must be set to a value greater than 0")
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

func (s *Simulator) Run() error {
	random := rand.New(rand.NewSource(s.seed))

	inputChan := make(chan ShuffleInput, s.numWorkers)
	resultChan := make(chan result.ShuffleResult, s.numWorkers)

	// Start consumer workers
	for range s.numWorkers {
		go PlayShuffleWorker(inputChan, resultChan)
	}

	shuffleId := uint(0)

	startTime := time.Now()

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
			return fmt.Errorf("error in shuffle %d: %w", shuffleResult.ShuffleId, shuffleResult.Error)
		}

		if s.verbose {
			log.Printf("Shuffle %d: Played %d rounds with final balance of: %d\n", shuffleResult.ShuffleId, shuffleResult.NumRounds, shuffleResult.Balance)
		}

		shuffleResults[shuffleResult.ShuffleId] = shuffleResult

		finish := false

		if shuffleResult.ShuffleId == countedShuffles {
			for i := countedShuffles; i < uint(len(shuffleResults)); i++ {
				if shuffleResults[i].NumRounds <= 0 {
					break
				}

				countedShuffles++

				if s.numShuffles > 0 {
					count++
					if count >= s.numShuffles {
						finish = true
					}
				}

				if s.numRounds > 0 {
					count += uint(shuffleResults[i].NumRounds)
					if count >= s.numRounds {
						finish = true
					}
				}

				if s.numHands > 0 {
					count += uint(shuffleResults[i].NumHands)
					if count >= s.numHands {
						finish = true
					}
				}

				if finish {
					log.Printf("Finished %d shuffles using %.3f seconds\n", countedShuffles, time.Since(startTime).Seconds())
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

	log.Printf("Average balance: %.2f\n", averageBalance)
	log.Printf("Total balance: %d\n", balanceSum)

	if s.csvFile != "" {
		csvExporter := exporter.NewCSVExporter(s.csvFile)
		log.Printf("Exporting results to CSV...\n")
		if err := csvExporter.Export(shuffleResults); err != nil {
			return fmt.Errorf("error exporting results to CSV: %w", err)
		}
	}

	log.Printf("Simulation completed successfully\n")
	return nil
}

func (s *Simulator) sendInput(inputChan chan<- ShuffleInput, shuffleId uint, random *rand.Rand) {
	player := person.NewPlayer(s.strategy)
	dealer := person.NewDealer()
	rules := NewRules(s.doubleAfterSplit, s.hitAfterSplitAce, s.splitAfterSplitAce, s.doubleAfterSplitAce, s.maxNumHands)
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
