package room

import (
	"context"
	"encoding/json"

	"github.com/yashikota/scene-hunter-backend/internal/util"
	"github.com/yashikota/scene-hunter-backend/model"
)

var ctx = context.Background()
var client = util.ConnectToUpstashRedis()

func CreateRoom(roomID string, user model.User) error {
	gameMaster, err := json.Marshal(user)
	if err != nil {
		return err
	}

	err = client.HSet(ctx, roomID, map[string]interface{}{
		"game_master":   gameMaster,
		"total_players": 0,
		"players":       "[]",
	}).Err()
	if err != nil {
		return err
	}

	return nil
}

func JoinRoom(roomID string, user model.User) error {
	// Convert user data to JSON
	playerData, err := json.Marshal(user)
	if err != nil {
		return err
	}

	// Get the current player list
	playerListData, err := client.HGet(ctx, roomID, "players").Result()
	if err != nil {
		return err
	}

	var players []json.RawMessage
	err = json.Unmarshal([]byte(playerListData), &players)
	if err != nil {
		return err
	}

	// Add the new player to the player list
	players = append(players, playerData)

	// Update the player list
	updatedPlayerListData, err := json.Marshal(players)
	if err != nil {
		return err
	}

	err = client.HSet(ctx, roomID, map[string]interface{}{
		"players":       updatedPlayerListData,
		"total_players": len(players),
	}).Err()
	if err != nil {
		return err
	}

	return nil
}

func CheckExistRoom(roomID string) bool {
	exists, err := client.Exists(ctx, roomID).Result()
	if err != nil || exists == 0 {
		return false
	}

	return true
}

func GetRoomUsers(roomID string) (*model.Room, error) {
	roomData, err := client.HGetAll(ctx, roomID).Result()
	if err != nil {
		return nil, err
	}

	// Convert the room data to a Room struct
	gameMaster := model.User{}
	err = json.Unmarshal([]byte(roomData["game_master"]), &gameMaster)
	if err != nil {
		return nil, err
	}

	var players []model.User
	err = json.Unmarshal([]byte(roomData["players"]), &players)
	if err != nil {
		return nil, err
	}

	room := &model.Room{
		GameMaster:   gameMaster,
		Players:      players,
		TotalPlayers: roomData["total_players"],
	}

	return room, nil
}

func DeleteRoom(roomID string) error {
	result := client.Del(ctx, roomID)

	return result.Err()
}
