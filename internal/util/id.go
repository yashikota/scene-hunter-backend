package util

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

var ctx, client = SetUpRedisClient()

func GenerateUserID(ttl int) (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	// Set TTL
	err = setUserID(id.String(), ttl)
	if err != nil {
		return "", err
	}

	return id.String(), nil
}

func setUserID(userID string, ttl int) error {
	now := time.Now().Unix()
	expire := now + int64(ttl)

	_, err := client.HSet(ctx, "UserID", userID, expire).Result()
	if err != nil {
		return err
	}

	return nil
}

func ExistUserID(userID string) (bool, error) {
	result, err := client.HExists(ctx, "UserID", userID).Result()
	if err != nil {
		return false, err
	}

	// Check TTL
	if result {
		expire, err := client.HGet(ctx, "UserID", userID).Int64()
		if err != nil {
			return false, err
		}

		now := time.Now().Unix()
		if now > expire {
			result = false
			client.HDel(ctx, "UserID", userID)
		}
	}

	return result, nil
}
