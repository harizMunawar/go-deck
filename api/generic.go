package api

import (
	"deck/card"
	"deck/database"
	"deck/utils"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const API_URL = "http://localhost:8000"

func CreateDeck(c *gin.Context) {
	deck := utils.CreateDeck()
	q := c.Request.URL.Query()

	if q.Get("shuffle") == "true" {
		deck = utils.Shuffle(deck)
	}

	r := rand.New(rand.NewSource(time.Now().Unix()))
	deckId := r.Intn(100000)

	deckDB := database.Deck{
		ID:         deckId,
		CardTotal:  len(deck),
		JokerTotal: 0,
	}

	database.DB.Create(deckDB)
	var cards = []database.Card{}
	for index, element := range deck {
		cards = append(cards, database.Card{
			Deck: database.Deck{
				ID:         deckDB.ID,
				CardTotal:  deckDB.CardTotal,
				JokerTotal: deckDB.JokerTotal,
			},
			Suit:     int(element.Suit),
			Rank:     int(element.Rank),
			Position: int(index + 1),
			Drawed:   false,
			Verbose:  element.Stringify(),
		})
	}
	database.DB.Create(&cards)

	response := gin.H{
		"message": "Deck created successfully",
		"deckId":  deckId,
	}

	c.JSON(http.StatusCreated, response)
}

func GetDeck(c *gin.Context) {
	var (
		deck        database.Deck
		card        database.Card
		cardsInDeck []database.Card
	)

	id, _ := strconv.Atoi(c.Param("id"))

	if query := database.DB.Model(database.Deck{}).First(&deck, id); query.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Can't find deck with that ID"})
		return
	}

	cards, _ := database.DB.Model(database.Card{}).Where("deck_id = ?", id).Rows()
	defer cards.Close()

	for cards.Next() {
		if database.DB.ScanRows(cards, &card) == nil {
			cardsInDeck = append(cardsInDeck, card)
		}
	}

	response := gin.H{
		"deck":  deck,
		"cards": cardsInDeck,
	}

	c.JSON(http.StatusOK, response)
}

func DeleteDeck(c *gin.Context) {
	var (
		deck database.Deck
	)
	id, _ := strconv.Atoi(c.Param("id"))

	if query := database.DB.Model(database.Deck{}).First(&deck, id); query.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Can't find deck with that ID"})
		return
	}

	database.DB.Delete(&deck, id)
	database.DB.Delete(database.Card{}, "deck_id = ?", id)
	database.DB.Delete(database.Blackjack{}, "deck_id = ?", id)

	response := gin.H{
		"message": "Deck deleted successfully",
	}
	c.JSON(http.StatusOK, response)
}

func ShuffleDeck(c *gin.Context) {
	var (
		deck            database.Deck
		newDeck         []card.Card
		cardPlaceholder database.Card
	)
	id, _ := strconv.Atoi(c.Param("id"))

	if query := database.DB.Model(database.Deck{}).First(&deck, id); query.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Can't find deck with that ID"})
		return
	}

	cards, _ := database.DB.Model(database.Card{}).Where("deck_id = ?", id).Rows()
	defer cards.Close()

	for cards.Next() {
		if database.DB.ScanRows(cards, &cardPlaceholder) == nil && !cardPlaceholder.Drawed {
			newDeck = append(newDeck, card.Card{
				Suit: card.Suit(cardPlaceholder.Suit),
				Rank: card.Rank(cardPlaceholder.Rank),
			})
		}
	}

	newDeck = utils.Shuffle(newDeck)
	for index, element := range newDeck {
		database.DB.Model(database.Card{}).Where("suit = ? AND rank = ? AND drawed = false", element.Suit, element.Rank).Update("position", index+1)
	}

	response := gin.H{
		"message": "Deck successfully shuffled",
	}
	c.JSON(http.StatusOK, response)
}

func DrawDeck(c *gin.Context) {
	var (
		deck      database.Deck
		minPos    int
		deckTotal int
		topDeck   database.Card
		cardsLeft int
	)
	id, _ := strconv.Atoi(c.Param("id"))

	if query := database.DB.Model(database.Deck{}).First(&deck, id); query.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Can't find deck with that ID"})
		return
	}

	if database.DB.Model(database.Card{}).Select("COUNT(drawed)").Where("drawed = ?", false).Scan(&cardsLeft); cardsLeft > 0 {
		database.DB.Model(database.Card{}).Select("MIN(position) AS topDeck").Where("deck_id = ? AND drawed = ?", id, false).Row().Scan(&minPos)
		database.DB.Model(database.Card{}).Where("deck_id = ? AND position = ?", deck.ID, minPos).Update("drawed", true).Scan(&topDeck)

		database.DB.Model(database.Card{}).Select("COUNT(id)").Where("deck_id = ? AND drawed = ?", id, false).Scan(&deckTotal)
		database.DB.Model(database.Deck{}).Where("id = ?", deck.ID).Update("card_total", deckTotal)

		c.JSON(http.StatusOK, gin.H{"card": topDeck})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"message": "This deck is empty"})
}
