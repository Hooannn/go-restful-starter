package main

import (
	"fmt"
	"log"

	"github.com/Hooannn/EventPlatform/configs"
	"github.com/Hooannn/EventPlatform/internal/entity"
	"github.com/Hooannn/EventPlatform/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	db, err := entity.InitDB()
	cfg := configs.LoadConfig(".env")
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	} else {
		log.Println("✅ Connected to database")
	}

	router := gin.Default()

	routes.RegisterRoutes(router, db)

	router.Run(fmt.Sprintf(":%s", cfg.Port))
}
