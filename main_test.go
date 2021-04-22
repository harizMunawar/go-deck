package main

import (
	"deck/card"
	"fmt"
)

func ExampleCard() {
	fmt.Println(card.Card{Suit: card.Spade, Rank: card.King}.Stringify())

	// Output:
	// King Of Spade
}
