package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rwiteshbera/MoneyTracker/middlewares"
	routes "github.com/rwiteshbera/MoneyTracker/routes"
)

func main() {
	err := godotenv.Load()
	checkError(err)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "6001"
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(middlewares.CorsMiddleware())
	router.Use(gin.Recovery())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	if err := router.Run(":" + PORT); err != nil {
		log.Panic("Failed to start server.")
	}
	fmt.Println("Server is listening!")
}

func checkError(err error) {
	if err != nil {
		log.Panic(err)
	}
}
