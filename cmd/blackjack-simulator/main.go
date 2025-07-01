package main

import (
	"log"

	"github.com/jljl1337/blackjack-simulator/internal/simulation"
)

func main() {
	s, err := simulation.NewSimulator()
	if err != nil {
		log.Fatalf("Error creating simulator: %v", err)
	}

	if err := s.Run(); err != nil {
		log.Fatalf("Error running simulation: %v", err)
	}
}
