package main

import (
	"deck/utils"
	"fmt"
)

func main() {
	deck := utils.CreateDeck(utils.Joker, utils.Shuffle)

	fmt.Println(deck)
	fmt.Println("DOROW")
	fmt.Println(utils.Draw(&deck))
	fmt.Println(deck)
}
