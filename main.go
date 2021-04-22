package main

import (
	"deck/deck"
	"fmt"
)

func main() {
	fmt.Println("Hello World")

	fmt.Println(deck.Card{Suit: deck.Heart, Rank: deck.Queen}.Stringify())
}
