package deck

import "fmt"

func ExampleCard() {
	fmt.Println(Card{Suit: Spade, Rank: Ace}.Stringify())

	// Output:
	// Ace Of Spade
}
