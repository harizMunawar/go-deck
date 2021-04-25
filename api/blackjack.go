package api

import (
	"deck/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func StartBlackjack(c *gin.Context) {
	var (
		card         database.Card
		player_hands []database.Card
		dealer_hands []database.Card
		drawedId     []int
	)

	deckId, _ := strconv.Atoi(c.Param("deckId"))

	game := database.BlackJack{
		DeckID: deckId,
		State:  false,
	}
	database.DB.Omit("PlayerHand", "DealerHand").Create(&game)

	rows, err := database.DB.Model(database.Card{}).Where("deck_id = ? AND drawed = false", deckId).Limit(4).Order("id ASC").Rows()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "An error happened, sorry :(",
		})
		return
	}

	i := 0
	for rows.Next() {
		if database.DB.ScanRows(rows, &card) == nil {
			drawedId = append(drawedId, card.ID)
			if i%2 == 0 {
				player_hands = append(player_hands, card)
			} else {
				dealer_hands = append(dealer_hands, card)
			}
		}
		i++
	}
	if i == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Cant initialize a game with that deck id. Either the deck is empty or cant find deck with that id",
		})
		return
	}
	database.DB.Model(&game).Association("PlayerHand").Append(player_hands)
	database.DB.Model(&game).Association("DealerHand").Append(dealer_hands)

	response := gin.H{
		"message": "Blackjack Game Started",
		"gameId":  game.ID,
	}
	c.JSON(http.StatusCreated, response)

	defer setDrawn(drawedId)
	defer rows.Close()
}

func setDrawn(listId []int) {
	for _, element := range listId {
		database.DB.Model(database.Card{}).Where("id = ?", element).Update("drawed", true)
	}
}
