package main

import (
	"log"
	"os"

	"ticket-system/database"
	"ticket-system/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load variables from .env file (useful for local development).
	_ = godotenv.Load()

	err := database.InitDB()
	if err != nil {
		log.Fatal("failed to initialize database: ", err)
	}

	r := gin.Default()
	routes.RegisterRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err = r.Run(":" + port)
	if err != nil {
		log.Fatal("failed to start server: ", err)
	}
}

