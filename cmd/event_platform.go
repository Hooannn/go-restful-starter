package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Hooannn/go-restful-starter/configs"
	"github.com/Hooannn/go-restful-starter/internal/entity"
	"github.com/Hooannn/go-restful-starter/internal/factory"
	"github.com/Hooannn/go-restful-starter/internal/redis"
	"github.com/Hooannn/go-restful-starter/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := configs.LoadConfig()

	if cfg.AppEnv == "development" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
		gin.DisableConsoleColor()
		f, _ := os.Create("server.log")
		gin.DefaultWriter = io.MultiWriter(f)
	}

	db, err := entity.InitDB()
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	} else {
		log.Println("✅ Connected to database")
	}

	redisClient, err := redis.InitRedis()
	if err != nil {
		log.Fatalf("❌ Failed to connect to redis: %v", err)
	} else {
		log.Println("✅ Connected to redis")
	}

	router := gin.Default()

	f := factory.NewFactory(db, redisClient)

	routes.RegisterRoutes(router, f)

	router.Run(fmt.Sprintf(":%s", cfg.Port))
}
