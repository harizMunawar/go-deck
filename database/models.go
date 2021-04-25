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

type BlackJack struct {
	ID          int    `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	DeckID      int    `json:"-"`
	Deck        Deck   `json:"deck" gorm:"foreignKey:DeckID;constraint:OnDelete:CASCADE;references:ID"`
	PlayerHand  []Card `json:"playerHand" gorm:"many2many:player_hands"`
	DealerHand  []Card `json:"dealerHand" gorm:"many2many:dealer_hands"`
	PlayerScore int    `json:"playerScore"`
	DealerScore int    `json:"dealerScore"`
	State       bool   `json:"state"`
}
