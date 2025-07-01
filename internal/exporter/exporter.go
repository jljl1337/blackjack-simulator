package exporter

import "github.com/jljl1337/blackjack-simulator/internal/result"

// Exporter is an interface for exporting simulation results.
type Exporter interface {
	Export(results []result.ShuffleResult) error
}
