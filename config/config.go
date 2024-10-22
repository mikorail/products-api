package config

import (
	"context"
	"log"
	"os"
	"products-api/global"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var RedisClient *redis.Client
var Ctx = context.Background()

// getEnv retrieves an environment variable or returns a default value if not set
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}

// Initialize database
func InitDB() {
	// Get the values from environment variables or use default values if not set
	host := getEnv(global.DB_URL, "localhost")
	user := getEnv(global.DB_USER, "user-message-api")
	password := getEnv(global.DB_PASSWORD, "th3password")
	dbname := getEnv(global.DB_NAME, "user-message-api")
	port := getEnv(global.DB_PORT, "5432")
	sslmode := getEnv(global.DB_SSL, "disable")

	// Construct the DSN (Data Source Name)
	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port + " sslmode=" + sslmode

	// Open the database connection
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Database connected")
}

// Initialize Redis
func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv(global.REDIS_URL),
		Password: os.Getenv(global.REDIS_PASSWORD),
		DB:       0,
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	log.Println("Redis connected")
}

// Init all configurations
func InitConfig() {
	InitDB()
	InitRedis()
}
