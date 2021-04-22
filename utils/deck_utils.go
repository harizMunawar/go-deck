package utils

import (
	"deck/card"
	"math/rand"
	"sort"
	"time"
)

func CreateDeck(options ...func([]card.Card) []card.Card) []card.Card {
	var cards []card.Card

	for _, suit := range [...]card.Suit{card.Spade, card.Diamond, card.Clover, card.Heart} {
		for rank := card.Ace; rank <= card.King; rank++ {
			cards = append(cards, card.Card{Suit: suit, Rank: rank})
		}
	}

	for _, opt := range options {
		cards = opt(cards)
	}

	return cards
}

func Shuffle(cards []card.Card) []card.Card {
	newDeck := make([]card.Card, len(cards))
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for i, j := range r.Perm(len(newDeck)) {
		newDeck[i] = cards[j]
	}

	return newDeck
}

func Sort(cards []card.Card) []card.Card {
	sort.Slice(cards, Less(cards))
	return cards
}

func Less(cards []card.Card) func(i, j int) bool {
	return func(i, j int) bool {
		return cardScore(cards[i]) < cardScore(cards[j])
	}
}

func Draw(cards *[]card.Card) card.Card {
	drawedCard := (*cards)[0]
	*cards = (*cards)[1:]
	return drawedCard
}

func cardScore(c card.Card) int {
	return int(c.Suit)*int(card.King) + int(c.Rank)
}
