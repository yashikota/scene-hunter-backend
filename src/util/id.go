package util

import (
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
)

var ctx, client = SetUpRedisClient()

func GenerateUserID(ttl int) (string, error) {
	id := ulid.Make()

	// Set TTL
	err := setUserID(id.String(), ttl)
	if err != nil {
		return "", fmt.Errorf("failed to set TTL for user id")
	}

	return id.String(), nil
}

func setUserID(userID string, ttl int) error {
	key := fmt.Sprintf("UserID:%s", userID)

	_, err := client.Set(ctx, key, true, time.Duration(ttl)*time.Second).Result()
	if err != nil {
		return fmt.Errorf("failed to generate user id")
	}

	return nil
}

func ExistUserID(userID string) (bool, error) {
	key := fmt.Sprintf("UserID:%s", userID)
	result, err := client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check user id")
	}

	if result == 0 {
		return false, nil
	}

	return true, nil
}
