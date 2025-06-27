package blackjack

import (
	_ "embed"
	"encoding/csv"
	"errors"
	"strings"

	"github.com/jljl1337/blackjack-simulator/internal/core"
)

//go:embed s17.csv
var s17csv string

type BasicStrategy struct {
	strategyTable map[string]map[string][]Action
}

// Create a new BasicStrategy instance with a predefined s17 strategy table.
func NewBasicStrategyS17() (*BasicStrategy, error) {
	return NewBasicStrategyFromCSV(s17csv)
}

// NewBasicStrategyFromCSV creates a new BasicStrategy instance from a CSV string.
func NewBasicStrategyFromCSV(csvString string) (*BasicStrategy, error) {
	// Create a CSV reader from the embedded string
	reader := csv.NewReader(strings.NewReader(csvString))

	// Read all records at once
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Convert records to a map of maps
	strategyMap, err := recordsToMapOfMaps(records)
	if err != nil {
		return nil, err
	}

	return &BasicStrategy{
		strategyTable: strategyMap,
	}, nil
}

// recordsToMapOfMaps converts a 2D slice of strings (CSV records) into a map
// of maps of Actions.
func recordsToMapOfMaps(records [][]string) (map[string]map[string][]Action, error) {
	if len(records) != 39 {
		return nil, errors.New("expected 39 rows in the CSV data")
	}

	result := make(map[string]map[string][]Action)
	headers := records[0]

	if len(headers) != 11 {
		return nil, errors.New("expected 11 columns in the CSV header")
	}

	// Loop through all rows starting from index 1 to skip header row
	for i := 1; i < len(records); i++ {
		record := records[i]
		if len(record) != 11 {
			return nil, errors.New("expected 11 columns in each data row")
		}

		rowMap := make(map[string][]Action)

		// Loop through each column in the row, starting from index 1
		for j := 1; j < len(headers); j++ {
			var actions, err = StringToActions(record[j])
			if err != nil {
				return nil, err
			}

			// Use the header as the key
			rowMap[headers[j]] = actions
		}

		result[record[0]] = rowMap
	}

	return result, nil
}

func (bs BasicStrategy) GetActions(
	playerHand core.Hand,
	dealerUpCard core.Card,
) ([]Action, error) {
	var playerHandValueKey string
	var err error

	if playerHand.IsPair() {
		playerHandValueKey, err = playerHand.PairString()
		if err != nil {
			return nil, err
		}
	} else {
		playerHandValueKey = playerHand.ValueString()
	}

	return bs.strategyTable[playerHandValueKey][dealerUpCard.ValueString()], nil
}
