package room

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/yashikota/scene-hunter-backend/internal/util"
	"github.com/yashikota/scene-hunter-backend/model"
)

var ctx, client = util.SetUpRedisClient()

func CreateRoom(roomID string, user model.User) error {
	newRoom := model.Room{
		GameMasterID: user.ID,
		TotalPlayers: 0,
		Users: map[string]model.User{
			user.ID: {
				ID:     user.ID,
				Name:   user.Name,
				Lang:   user.Lang,
				Status: "active",
				Score:  []float32{},
				Photo:  []string{},
			},
		},
	}

	roomJSON, err := json.Marshal(newRoom)
	if err != nil {
		return err
	}

	err = client.JSONSet(ctx, roomID, "$", string(roomJSON)).Err()
	if err != nil {
		return err
	}

	return nil
}

func JoinRoom(roomID string, user model.User) error {
	newPlayer := model.User{
		ID:     user.ID,
		Name:   user.Name,
		Lang:   user.Lang,
		Status: "active",
		Score:  []float32{},
		Photo:  []string{},
	}

	playerJSON, err := json.Marshal(newPlayer)
	if err != nil {
		return err
	}

	err = client.JSONSet(ctx, roomID, "$.users."+user.ID, string(playerJSON)).Err()
	if err != nil {
		log.Printf("Failed to join the room: %v", err)
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
	result, err := client.JSONGet(ctx, roomID, fmt.Sprintf("$.users.%s", userID)).Result()
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

func ChangeGameMaster(roomID string, userID string) error {
	err := client.JSONSet(ctx, roomID, "$.game_master_id", fmt.Sprintf("\"%s\"", userID)).Err()
	if err != nil {
		return err
	}

	return nil
}

func DeleteRoomUser(roomID string, userID string) error {
	err := client.JSONDel(ctx, roomID, fmt.Sprintf("$.users.%s", userID)).Err()
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
