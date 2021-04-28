package api

import (
	"deck/database"
	"deck/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func InitBlackjack(c *gin.Context) {
	var (
		deck       database.Deck
		deckJSON   map[string]interface{}
		playerHand []database.Card
		dealerHand []database.Card
	)

	deckRequest, err := http.Get(utils.BASE_URL + "/deck/create?shuffle=true")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Oops, something wrong happened!",
		})
		return
	}
	utils.ReadToMap(deckRequest, &deckJSON)
	deckId := int(deckJSON["deckId"].(float64))
	database.DB.Model(database.Deck{}).Where("id = ?", deckId).Scan(&deck)

	for i := 0; i < 2; i++ {
		playerHand, _ = draw(deck, playerHand)
		dealerHand, _ = draw(deck, dealerHand)
	}

	playerScore, _ := score_counter(playerHand)
	dealerScore, _ := score_counter(dealerHand)

	game := database.Blackjack{
		Deck:        deck,
		DeckID:      deck.ID,
		PlayerHand:  playerHand,
		PlayerScore: playerScore,
		DealerHand:  dealerHand,
		DealerScore: dealerScore,
	}

	database.DB.Create(&game)

	c.JSON(http.StatusCreated, gin.H{
		"blackjack": game,
	})
}

func BlackjackStatus(c *gin.Context) {
	var game database.Blackjack
	forceFinished := false

	id, _ := strconv.Atoi(c.Param("gameid"))

	q := c.Request.URL.Query()
	if q.Get("end") == "true" {
		forceFinished = true
	}

	if query := database.DB.Model(database.Blackjack{}).Preload(clause.Associations).Where("id = ?", id).Find(&game); query.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Can't find blackjack game with that ID",
		})
		return
	}
	playerScore, _ := score_counter(game.PlayerHand)
	dealerScore, _ := score_counter(game.DealerHand)

	winner, message := check_status(playerScore, dealerScore)
	finished := check_finished(playerScore, dealerScore, forceFinished)

	database.DB.Model(&game).Updates(database.Blackjack{
		Winner:      winner,
		Message:     message,
		Finished:    finished,
		PlayerScore: playerScore,
		DealerScore: dealerScore,
	})

	c.JSON(http.StatusOK, gin.H{
		"blackjack": game,
	})
}

func BlackjackHit(c *gin.Context) {
	var (
		game database.Blackjack
		deck database.Deck
	)

	gameId, _ := strconv.Atoi(c.Param("gameid"))
	query := database.DB.Model(database.Blackjack{}).Preload(clause.Associations).Where("id = ?", gameId).Scan(&game)
	if query.Error != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Can't found blackjack game with that ID",
		})
	}

	database.DB.Model(database.Deck{}).Where("id = ?", game.DeckID).Scan(&deck)
	playerHand, _ := draw(deck, game.PlayerHand)
	database.DB.Model(&game).Association("PlayerHand").Append(playerHand)

	c.JSON(http.StatusOK, gin.H{
		"message": "Hit success",
	})
}

func DealerAI(c *gin.Context) {
	var (
		game database.Blackjack
		deck database.Deck
	)

	forceFinished := false

	q := c.Request.URL.Query()
	if q.Get("stand") == "true" {
		forceFinished = true
	}

	gameId, _ := strconv.Atoi(c.Param("gameid"))
	get_data(&game, &deck, gameId, c)

	dealerScore, containAces := score_counter(game.DealerHand)

	for dealerScore <= 16 || (dealerScore <= 17 && containAces) {
		dealerHand, _ := draw(deck, game.DealerHand)
		database.DB.Model(&game).Association("DealerHand").Replace(dealerHand)
		get_data(&game, &deck, gameId, c)
		dealerScore, containAces = score_counter(game.DealerHand)
	}

	if forceFinished {
		database.DB.Model(&game).Update("finished", true)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Dealer have made their move",
	})
}

func score_counter(hand []database.Card) (int, bool) {
	var (
		score int
	)

	rank := []int{0, 11, 2, 3, 4, 5, 6, 7, 8, 9, 10, 10, 10, 10}
	containsAce := false

	for _, card := range hand {
		score += rank[card.Rank]
		if card.Rank == 1 {
			containsAce = true
		}
	}

	if containsAce && score > 21 {
		score -= 10
	}

	return score, containsAce
}

func draw(deck database.Deck, hand []database.Card) ([]database.Card, error) {
	var maps map[string]interface{}

	resp, err := http.Get(utils.BASE_URL + "/deck/" + strconv.Itoa(deck.ID) + "/draw")
	if err != nil {
		return hand, err
	}
	utils.ReadToMap(resp, &maps)
	maps = maps["card"].(map[string]interface{})

	card := database.Card{
		ID:       int(maps["id"].(float64)),
		DeckID:   deck.ID,
		Deck:     deck,
		Suit:     int(maps["suit"].(float64)),
		Rank:     int(maps["rank"].(float64)),
		Position: int(maps["position"].(float64)),
		Verbose:  maps["verbose"].(string),
	}
	hand = append(hand, card)

	return hand, nil
}

func check_status(playerScore int, dealerScore int) (string, string) {
	if playerScore > 21 {
		return "Dealer", "Player Busted"
	}
	if dealerScore > 21 {
		return "Player", "Dealer Busted"
	}
	if dealerScore == 21 {
		return "Dealer", "Dealer Mendapat Blackjack"
	}
	if playerScore == 21 {
		return "Player", "Player Mendapat Blackjack"
	}
	if playerScore > dealerScore {
		return "Player", "Skor Player Lebih Besar Dari Dealer"
	}
	if dealerScore > playerScore {
		return "Dealer", "Skor Dealer Lebih Besar Dari Player"
	}
	if dealerScore == playerScore {
		return "Draw", "Skor Player & Dealer Sama"
	}

	return "", ""
}

func check_finished(playerScore int, dealerScore int, forceFinished bool) bool {
	if dealerScore == 21 || playerScore == 21 || playerScore >= 21 || (dealerScore >= 21 && forceFinished) || forceFinished {
		return true
	}
	return false
}

func get_data(game *database.Blackjack, deck *database.Deck, gameId int, c *gin.Context) {
	query := database.DB.Model(database.Blackjack{}).Preload(clause.Associations).Where("id = ?", gameId).Find(&game)
	if query.Error != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Can't found blackjack game with that ID",
		})
	}

	database.DB.Model(database.Deck{}).Where("id = ?", game.DeckID).Scan(&deck)
}
