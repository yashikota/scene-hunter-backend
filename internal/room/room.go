package room

import (
	"encoding/json"
	"fmt"

	"github.com/yashikota/scene-hunter-backend/internal/util"
	"github.com/yashikota/scene-hunter-backend/model"
)

var ctx, client = util.SetUpRedisClient()

func CreateRoom(roomID string, user model.User) error {
	newRoom := fmt.Sprintf(`{"game_master":{"id":"%s","name":"%s","lang":"%s"},"players":[],"total_players":0}`, user.ID, user.Name, user.Lang)

	err := client.JSONSet(ctx, roomID, "$", newRoom).Err()
	if err != nil {
		return err
	}

	return nil
}

func JoinRoom(roomID string, user model.User) error {
	newPlayer := fmt.Sprintf(`{"id":"%s","name":"%s","lang":"%s"}`, user.ID, user.Name, user.Lang)

	err := client.JSONArrAppend(ctx, roomID, "$.players", newPlayer).Err()
	if err != nil {
		return err
	}

	err = client.JSONNumIncrBy(ctx, roomID, "$.total_players", 1).Err()
	if err != nil {
		return err
	}

	return nil
}

func CheckExistRoom(roomID string) bool {
	result := client.Exists(ctx, roomID)

	return result.Val() == 1
}

func CheckExistUser(roomID string, userID string) (bool, error) {
	result, err := client.JSONGet(ctx, roomID, fmt.Sprintf("$.players[?(@.id=='%s')]", userID)).Result()
	if err != nil {
		return false, err
	}

	if result == "" {
		return false, nil
	}

	return true, nil
}

func GetRoomUsers(roomID string) ([]*model.Room, error) {
	jsonData, err := client.JSONGet(ctx, roomID, "$").Result()
	if err != nil {
		return nil, err
	}

	var roomData []*model.Room
	err = json.Unmarshal([]byte(jsonData), &roomData)
	if err != nil {
		return nil, err
	}

	return roomData, nil
}

func DeleteRoomUser(roomID string, userID string) error {
	err := client.JSONArrPop(ctx, roomID, "$.players", 0).Err()
	if err != nil {
		return err
	}

	err = client.JSONNumIncrBy(ctx, roomID, "$.total_players", -1).Err()
	if err != nil {
		return err
	}

	return nil
}

func DeleteRoom(roomID string) error {
	err := client.Del(ctx, roomID).Err()
	if err != nil {
		return err
	}

	err = util.DeleteDir(fmt.Sprintf("uploads/%s", roomID))
	if err != nil {
		return err
	}

	return nil
}
