package exporter

import (
	"encoding/csv"
	"os"
	"strconv"

	"github.com/jljl1337/blackjack-simulator/internal/result"
)

type CSVExporter struct {
	filePath string
}

func NewCSVExporter(filePath string) *CSVExporter {
	return &CSVExporter{
		filePath: filePath,
	}
}

func (e CSVExporter) Export(results []result.ShuffleResult) error {
	data := make([][]string, 0, len(results))

	data = append(data, []string{
		"shuffle_id",
		"round_id",
		"hand_id",
		"dealer_hand",
		"player_hand",
		"dealer_value",
		"player_value",
		"player_actions",
		"bet_placed",
		"bet",
	})

	for resultID, result := range results {
		for roundID, roundResult := range result.RoundResults {
			for handID, playerHand := range roundResult.PlayerHands {
				dealerHand := roundResult.DealerHand.String()
				playerHands := playerHand.Hand.String()
				dealerHandValue := strconv.Itoa(roundResult.DealerHand.Value())
				playerHandValue := strconv.Itoa(playerHand.Value())
				playerActions := playerHand.GetActions()
				betPlaced := strconv.Itoa(playerHand.GetBetPlaced())
				bet := strconv.Itoa(playerHand.GetBet())

				// Convert player actions to a string representation
				playerActionsString := ""

				for _, action := range playerActions {
					playerActionsString += action.String() + ";"
				}
				// Remove the trailing semicolon if there are any actions
				if len(playerActionsString) > 0 {
					playerActionsString = playerActionsString[:len(playerActionsString)-1]
				}

				data = append(data, []string{
					strconv.Itoa(resultID),
					strconv.Itoa(roundID),
					strconv.Itoa(handID),
					dealerHand,
					playerHands,
					dealerHandValue,
					playerHandValue,
					playerActionsString,
					betPlaced,
					bet,
				})
			}
		}
	}

	if err := e.saveToCSV(data); err != nil {
		return err
	}

	return nil
}

func (e CSVExporter) saveToCSV(data [][]string) error {
	// Create or truncate output CSV file
	file, err := os.Create(e.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Initialize CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write all rows at once
	if err := writer.WriteAll(data); err != nil {
		return err
	}

	return nil
}
