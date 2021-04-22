package card

type Suit uint8

const (
	Spade Suit = iota
	Diamond
	Clover
	Heart
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
	Joker
)

type Card struct {
	Suit
	Rank
}

func (c Card) Stringify() string {
	if c.Rank == Joker {
		return c.Rank.String()
	}

	return c.Rank.String() + " Of " + c.Suit.String()
}
