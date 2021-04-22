package card

type Suit uint8

const (
	Spade Suit = iota
	Diamond
	Clover
	Heart
	Joker
)

type Rank uint8

const (
	_ Rank = iota
	Ace
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

type Card struct {
	Suit
	Rank
}

func (c Card) Stringify() string {
	if c.Suit == Joker {
		return c.Suit.String()
	}

	return c.Rank.String() + " Of " + c.Suit.String()
}
