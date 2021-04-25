package database

type Deck struct {
	ID         int `json:"id" gorm:"unique;primaryKey"`
	CardTotal  int `json:"cardTotal"`
	JokerTotal int `json:"jokerTotal"`
}

type Card struct {
	ID       int    `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	DeckID   int    `json:"-"`
	Deck     Deck   `json:"-" gorm:"foreignKey:DeckID;constraint:OnDelete:CASCADE;references:ID"`
	Suit     int    `json:"suit"`
	Rank     int    `json:"rank"`
	Position int    `json:"position"`
	Drawed   bool   `json:"isCardDrawed"`
	Verbose  string `json:"verbose"`
}
