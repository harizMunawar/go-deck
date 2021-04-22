package main

import (
	"deck/card"
	"fmt"
)

func main() {
	fmt.Println("Hello World")

	fmt.Println(card.Card{Suit: card.Heart, Rank: card.Queen}.Stringify())
}
