package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Hooannn/go-restful-starter/configs"
	"github.com/Hooannn/go-restful-starter/internal/entity"
	"github.com/Hooannn/go-restful-starter/internal/factory"
	"github.com/Hooannn/go-restful-starter/internal/middleware"
	"github.com/Hooannn/go-restful-starter/internal/opentelemetry"
	"github.com/Hooannn/go-restful-starter/internal/redis"
	"github.com/Hooannn/go-restful-starter/internal/routes"
	"github.com/Hooannn/go-restful-starter/internal/worker"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
)

var (
	cfg    = configs.LoadConfig()
	tracer = otel.Tracer(cfg.AppName)
)

func main() {
	tp, err := opentelemetry.InitTracerProvider()
	if err != nil {
		log.Printf("❌ Failed to init tracer provider: %v", err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
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
	router.Use(middleware.WithTracer(tracer))

	f := factory.NewFactory(db, redisClient)

	routes.RegisterRoutes(router, f)

	go worker.Bootstrap(cfg, db)
	router.Run(fmt.Sprintf(":%s", cfg.Port))
}
