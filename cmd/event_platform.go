package main

import (
	"fmt"
	"log"

	"github.com/Hooannn/EventPlatform/configs"
	"github.com/Hooannn/EventPlatform/internal/entity"
	"github.com/Hooannn/EventPlatform/internal/factory"
	"github.com/Hooannn/EventPlatform/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := configs.LoadConfig(".env")
	db, err := entity.InitDB()
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	} else {
		log.Println("✅ Connected to database")
	}
	router := gin.Default()

	f := factory.NewFactory(db)

	routes.RegisterRoutes(router, f)

	router.Run(fmt.Sprintf(":%s", cfg.Port))
}
