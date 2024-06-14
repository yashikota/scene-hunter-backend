package util

import (
	"context"
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

func getUpstashSecret() (string, string) {
	loadEnv()

	UPSTASH_REDIS_URL := os.Getenv("UPSTASH_REDIS_URL")
	UPSTASH_REDIS_TOKEN := os.Getenv("UPSTASH_REDIS_TOKEN")

	return UPSTASH_REDIS_URL, UPSTASH_REDIS_TOKEN
}

func connectToUpstash() *redis.Client {
	upstashURL, upstashRedisToken := getUpstashSecret()
	upstashUrl := fmt.Sprintf("rediss://default:%s@%s", upstashRedisToken, upstashURL)
	opt, _ := redis.ParseURL(upstashUrl)
	client := redis.NewClient(opt)

	return client
}

func SetUpRedisClient() (context.Context, *redis.Client) {
	var ctx = context.Background()
	var client = connectToUpstash()

	return ctx, client
}
