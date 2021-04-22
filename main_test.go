package main

import (
	"deck/deck"
	"fmt"
)

func ExampleCard() {
	fmt.Println(deck.Card{Suit: deck.Spade, Rank: deck.King}.Stringify())

	// Output:
	// King Of Spade
}
