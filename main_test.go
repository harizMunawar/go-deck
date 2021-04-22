package main

import (
	"deck/card"
	"deck/utils"
	"fmt"
)

func ExampleCard() {
	fmt.Println(card.Card{Suit: card.Spade, Rank: card.King}.Stringify())

	// Output:
	// King Of Spade
}

func CreateDeck() {
	fmt.Println(len(utils.CreateDeck()))

	// Output:
	// 52
}

func DrawCard() {
	deck := utils.CreateDeck(utils.Shuffle)
	predictDraw := deck[0]
	cardDrawed := utils.Draw(&deck)

	fmt.Println((predictDraw == cardDrawed) && (deck[0] != predictDraw))
	// Output:
	// true
}
