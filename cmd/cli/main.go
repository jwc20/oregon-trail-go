package main

import (
	"fmt"
	"os"

	trail "github.com/jwc20/oregontrail"
)

type SimpleStore struct{}

func (s *SimpleStore) SaveState(state trail.GameState) error {
	return nil
}

func (s *SimpleStore) LoadState() (trail.GameState, error) {
	return trail.GameState{}, nil
}

func main() {
	fmt.Println("*************************************")
	fmt.Println("*     THE OREGON TRAIL GAME         *")
	fmt.Println("*************************************")
	fmt.Println()

	store := &SimpleStore{}
	game := trail.NewCLI(store, os.Stdin, os.Stdout)
	game.PlaySVT()
}
