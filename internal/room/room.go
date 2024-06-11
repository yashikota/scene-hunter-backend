package room

import (
	"context"

	"github.com/yashikota/scene-hunter-backend/internal/util"
	"github.com/yashikota/scene-hunter-backend/model"
)

var ctx = context.Background()
var client = util.ConnectToUpstashRedis()

func CreateRoom(roomID string) error {
	result := client.HSet(ctx, roomID, "", "")

	return result.Err()
}

func CheckExistRoom(roomID string) bool {
	result := client.HExists(ctx, roomID, "")

	return result.Val()
}

func JoinRoom(roomID string, user model.User) error {
	result := client.HSet(ctx, roomID, user.ID, user.Name)

	return result.Err()
}

func DeleteRoom(roomID string) error {
	result := client.Del(ctx, roomID)

	return result.Err()
}
