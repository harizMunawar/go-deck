package main

import (
	"deck/api"
	"deck/database"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Build compiled")
	database.InitDB()
	fmt.Println("Connected to database")
	router := gin.Default()

	go urls(router)
	go blackjack_urls(router)

	router.Run(":8000")
}

func urls(r *gin.Engine) {
	r.GET("/deck/create", api.CreateDeck)
	r.GET("/deck/:id", api.GetDeck)
	r.DELETE("/deck/:id", api.DeleteDeck)
	r.GET("/deck/:id/shuffle", api.ShuffleDeck)
	r.GET("/deck/:id/draw", api.DrawDeck)
}

func blackjack_urls(r *gin.Engine) {
	r.GET("/blackjack/create", api.InitBlackjack)
	r.GET("/blackjack/:gameid/status", api.BlackjackStatus)
	r.GET("/blackjack/:gameid/hit", api.BlackjackHit)
	r.GET("/blackjack/:gameid/dealer-ai", api.DealerAI)
}
