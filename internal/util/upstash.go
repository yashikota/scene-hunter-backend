package util

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func getUpstashRedisSecret() (string, string) {
	loadEnv()

	UPSTASH_REDIS_URL := os.Getenv("UPSTASH_REDIS_URL")
	UPSTASH_REDIS_TOKEN := os.Getenv("UPSTASH_REDIS_TOKEN")

	return UPSTASH_REDIS_URL, UPSTASH_REDIS_TOKEN
}

func ConnectToUpstashRedis() *redis.Client {
	upstashRedisURL, upstashRedisToken := getUpstashRedisSecret()
	upstashUrl := fmt.Sprintf("rediss://default:%s@%s", upstashRedisToken, upstashRedisURL)
	opt, _ := redis.ParseURL(upstashUrl)
	client := redis.NewClient(opt)

	return client
}
