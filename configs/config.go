package configs

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName                       string
	AppEnv                        string
	Port                          string
	DatabaseConnectionString      string
	ResetPasswordOTPExpireMinutes int
	JWTAccessTokenSecret          string
	JWTRefreshTokenSecret         string
	JWTAccessTokenExpireHours     int
	JWTRefreshTokenExpireHours    int
	RedisAddress                  string
	RedisUsername                 string
	RedisPassword                 string
	RedisDatabase                 int
	EmailSender                   string
	EmailPassword                 string
	SMTPHost                      string
	SMTPPort                      int
	DefaultCacheExpireMinutes     int
	JaegerHost                    string
	JaegerPort                    int
	RootUser                      string
	RootUserPassword              string
	RootRoleName                  string
	RootRoleDescription           string
	RootRolePermissions           string
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

func LoadConfig() *Config {
	once.Do(func() {
		if err := godotenv.Load(".env"); err != nil {
			log.Println("Error loading .env file, using default values")
		}
		cfg = &Config{
			AppName:                       getEnv("APP_NAME", "EventPlatform"),
			AppEnv:                        getEnv("APP_ENV", "development"),
			Port:                          getEnv("PORT", "8080"),
			ResetPasswordOTPExpireMinutes: getEnvAsInt("RESET_PASSWORD_OTP_EXPIRE_MINUTES", 5),
			DatabaseConnectionString:      getEnv("DATABASE_CONNECTION_STRING", "host=localhost port=5432 user=postgres password=postgres dbname=mydb sslmode=disable"),
			JWTAccessTokenSecret:          getEnv("JWT_ACCESS_TOKEN_SECRET", "defaultAccessTokenSecret"),
			JWTRefreshTokenSecret:         getEnv("JWT_REFRESH_TOKEN_SECRET", "defaultRefreshTokenSecret"),
			JWTAccessTokenExpireHours:     getEnvAsInt("JWT_ACCESS_TOKEN_EXPIRE_HOURS", 1),
			JWTRefreshTokenExpireHours:    getEnvAsInt("JWT_REFRESH_TOKEN_EXPIRE_HOURS", 24*7),
			RedisAddress:                  getEnv("REDIS_ADDRESS", "localhost:6379"),
			RedisDatabase:                 getEnvAsInt("REDIS_DB", 0),
			RedisUsername:                 getEnv("REDIS_USERNAME", "EventPlatform"),
			RedisPassword:                 getEnv("REDIS_PASSWORD", ""),
			EmailSender:                   getEnv("EMAIL_SENDER", ""),
			EmailPassword:                 getEnv("EMAIL_PASSWORD", ""),
			SMTPHost:                      getEnv("SMTP_HOST", "smtp.example.com"),
			SMTPPort:                      getEnvAsInt("SMTP_PORT", 587),
			DefaultCacheExpireMinutes:     getEnvAsInt("DEFAULT_CACHE_EXPIRE_MINUTES", 10),
			JaegerHost:                    getEnv("JAEGER_HOST", "localhost"),
			JaegerPort:                    getEnvAsInt("JAEGER_PORT", 4317),
			RootUser:                      getEnv("ROOT_USER", "root@gmail.com"),
			RootUserPassword:              getEnv("ROOT_USER_PASSWORD", "root"),
			RootRoleName:                  getEnv("ROOT_ROLE_NAME", "root"),
			RootRoleDescription:           getEnv("ROOT_ROLE_DESCRIPTION", "Root role of the system"),
			RootRolePermissions:           getEnv("ROOT_ROLE_PERMISSIONS", "read:*,create:*,delete:*,update:*"),
		}
	})
	return cfg
}
