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
	urls(router)
	router.Run(":8000")
}

func urls(r *gin.Engine) {
	r.GET("/create-deck", api.CreateDeck)
	r.GET("/deck/:id", api.GetDeck)
	r.DELETE("/deck/:id", api.DeleteDeck)
	r.GET("/deck/:id/shuffle", api.ShuffleDeck)
	r.GET("/deck/:id/draw", api.DrawDeck)
}
