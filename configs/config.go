package configs

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName                   string
	Port                      string
	DatabaseConnectionString  string
	JWTAccessTokenSecret      string
	JWTRefreshTokenSecret     string
	JWTAccessTokenExpireTime  int
	JWTRefreshTokenExpireTime int
	RedisHost                 string
	RedisPort                 string
	RedisPassword             string
}

var (
	cfg  *Config
	once sync.Once
)

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func LoadConfig(envFile string) *Config {
	once.Do(func() {
		if err := godotenv.Load(envFile); err != nil {
			log.Println("Error loading .env file, using default values")
		}
		cfg = &Config{
			AppName:                   getEnv("APP_NAME", "EventPlatform"),
			Port:                      getEnv("PORT", "8080"),
			DatabaseConnectionString:  getEnv("DATABASE_CONNECTION_STRING", "host=localhost port=5432 user=postgres password=postgres dbname=mydb sslmode=disable"),
			JWTAccessTokenSecret:      getEnv("JWT_ACCESS_TOKEN_SECRET", "defaultAccessTokenSecret"),
			JWTRefreshTokenSecret:     getEnv("JWT_REFRESH_TOKEN_SECRET", "defaultRefreshTokenSecret"),
			JWTAccessTokenExpireTime:  getEnvAsInt("JWT_ACCESS_TOKEN_EXPIRE_TIME", 3600),
			JWTRefreshTokenExpireTime: getEnvAsInt("JWT_REFRESH_TOKEN_EXPIRE_TIME", 86400),
			RedisHost:                 getEnv("REDIS_HOST", "localhost"),
			RedisPort:                 getEnv("REDIS_PORT", "6379"),
			RedisPassword:             getEnv("REDIS_PASSWORD", ""),
		}
	})
	return cfg
}
