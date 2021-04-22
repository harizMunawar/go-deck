package main

import (
	"deck/card"
	"fmt"
)

func main() {
	fmt.Println("Hello World")

	fmt.Println(card.Card{Suit: card.Heart, Rank: card.Queen}.Stringify())
	deck := CreateDeck()

	fmt.Println(deck)
	fmt.Println(len(deck))
}

func CreateDeck() []card.Card {
	var cards []card.Card

	for _, suit := range [...]card.Suit{card.Spade, card.Diamond, card.Clover, card.Heart} {
		for rank := card.Ace; rank <= card.King; rank++ {
			cards = append(cards, card.Card{Suit: suit, Rank: rank})
		}
	}

	return cards
}
