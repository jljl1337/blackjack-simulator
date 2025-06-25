package main

import (
	"fmt"
	"math/rand"

	"github.com/jljl1337/blackjack-simulator/internal/core"
)

func main() {
	fmt.Println("Hello, World!")

	shoe := core.NewShoe(8, 0.75, rand.New(rand.NewSource(1337)))
	fmt.Println("First card in shoe:", shoe.Deal())
}
