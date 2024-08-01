package room

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yashikota/scene-hunter-backend/model"
	"github.com/yashikota/scene-hunter-backend/src/util"
)

var ctx, client = util.SetUpRedisClient()

func CreateRoom(roomID string, user model.User) error {
	newRoom := model.Room{
		GameMasterID: user.ID,
		TotalPlayers: 0,
		GameRounds:   3,
		Users: map[string]model.User{
			user.ID: {
				ID:              user.ID,
				Name:            user.Name,
				Lang:            user.Lang,
				Status:          "active",
				PhotoScoreIndex: 0,
				Score:           map[int]float32{},
				Photo:           map[int]string{},
			},
		},
	}

	roomJSON, err := json.Marshal(newRoom)
	if err != nil {
		return fmt.Errorf("failed to create the room")
	}

	key := fmt.Sprintf("RoomID:%s", roomID)
	err = client.JSONSet(ctx, key, "$", string(roomJSON)).Err()
	if err != nil {
		return fmt.Errorf("failed to create the room")
	}

	// Set the expiration time to 24 hours
	err = client.Expire(ctx, roomID, 24 * time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to set the expiration time")
	}

	return nil
}

func JoinRoom(roomID string, user model.User) error {
	newPlayer := model.User{
		ID:     user.ID,
		Name:   user.Name,
		Lang:   user.Lang,
		Status: "active",
		Score:  map[int]float32{},
		Photo:  map[int]string{},
	}

	playerJSON, err := json.Marshal(newPlayer)
	if err != nil {
		return fmt.Errorf("failed to join the room")
	}

	key := fmt.Sprintf("RoomID:%s", roomID)
	err = client.JSONSet(ctx, key, "$.users."+user.ID, string(playerJSON)).Err()
	if err != nil {
		return fmt.Errorf("failed to join the room")
	}

	err = client.JSONNumIncrBy(ctx, key, "$.total_players", 1).Err()
	if err != nil {
		return fmt.Errorf("failed to increment the total players")
	}

	return nil
}

func CheckExistRoom(roomID string) (bool, error) {
	key := fmt.Sprintf("RoomID:%s", roomID)
	result := client.Exists(ctx, key)
	if result.Val() == 0 {
		return false, fmt.Errorf("room not found: %s", roomID)
	}

	return true, nil
}

func CheckExistUser(roomID string, userID string) (bool, error) {
	key := fmt.Sprintf("RoomID:%s", roomID)
	result, err := client.JSONGet(ctx, key, fmt.Sprintf("$.users.%s", userID)).Result()

	if err != nil {
		return false, err
	}
	if result == "[]" {
		return false, nil
	}

	return true, nil
}

func GetRoomUsers(roomID string) ([]*model.Room, error) {
	key := fmt.Sprintf("RoomID:%s", roomID)
	jsonData, err := client.JSONGet(ctx, key, "$").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get the room data")
	}

	var roomData []*model.Room
	err = json.Unmarshal([]byte(jsonData), &roomData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal the room data")
	}

	return roomData, nil
}

func CheckGameMaster(roomID string, userID string) (bool, int, error) {
	key := fmt.Sprintf("RoomID:%s", roomID)
	gameMasterID, err := client.JSONGet(ctx, key, "$.game_master_id").Result()
	if err != nil {
		return false, http.StatusInternalServerError, fmt.Errorf("failed to get the game master ID")
	}

	gameMasterID = gameMasterID[2 : len(gameMasterID)-2]
	if gameMasterID != userID {
		return false, http.StatusForbidden, fmt.Errorf("you are not the game master")
	}

	return true, http.StatusOK, nil
}

func ChangeGameMaster(roomID string, userID string) error {
	key := fmt.Sprintf("RoomID:%s", roomID)
	err := client.JSONSet(ctx, key, "$.game_master_id", fmt.Sprintf("\"%s\"", userID)).Err()
	if err != nil {
		return fmt.Errorf("failed to change the game master")
	}

	return nil
}

func UpdateRounds(roomID string, rounds string) error {
	if !util.ValidateRounds(rounds) {
		return fmt.Errorf("rounds must be between 1 and 10")
	}

	key := fmt.Sprintf("RoomID:%s", roomID)
	err := client.JSONSet(ctx, key, "$.game_rounds", rounds).Err()
	if err != nil {
		return fmt.Errorf("failed to update the rounds")
	}

	return nil
}

func AddRoomUserPhotoAndScore(roomID string, userID string, photoURL string) error {
	key := fmt.Sprintf("RoomID:%s", roomID)
	photoScoreIndex, err := client.JSONGet(ctx, key, fmt.Sprintf("$.users.%s.photo_score_index", userID)).Result()
	if err != nil {
		return fmt.Errorf("failed to get the player's score index")
	}

	// [N] to remove brackets
	photoScoreIndex = photoScoreIndex[1 : len(photoScoreIndex)-1]

	err = client.JSONSet(ctx, roomID, fmt.Sprintf("$.users.%s.photo.%s", userID, photoScoreIndex), fmt.Sprintf(
		"\"%s\"", photoURL)).Err()
	if err != nil {
		return fmt.Errorf("failed to set the player's photo")
	}

	score := util.GenerateRandomScore()
	err = client.JSONSet(ctx, roomID, fmt.Sprintf("$.users.%s.score.%s", userID, photoScoreIndex), fmt.Sprintf(
		"%.4f", score)).Err()
	if err != nil {
		return fmt.Errorf("failed to set the player's score")
	}

	err = client.JSONNumIncrBy(ctx, roomID, fmt.Sprintf("$.users.%s.photo_score_index", userID), 1).Err()
	if err != nil {
		return fmt.Errorf("failed to increment the player's score index")
	}

	return nil
}

func DeleteRoomUser(roomID string, userID string) error {
	key := fmt.Sprintf("RoomID:%s", roomID)
	err := client.JSONDel(ctx, key, fmt.Sprintf("$.users.%s", userID)).Err()
	if err != nil {
		return fmt.Errorf("failed to delete the user from the room")
	}

	err = client.JSONNumIncrBy(ctx, key, "$.total_players", -1).Err()
	if err != nil {
		return fmt.Errorf("failed to decrement the total players")
	}

	return nil
}

func DeleteRoom(roomID string) error {
	key := fmt.Sprintf("RoomID:%s", roomID)
	err := client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete the room")
	}

	err = util.DeleteDir(fmt.Sprintf("uploads/%s", roomID))
	if err != nil {
		return fmt.Errorf("failed to delete the room directory")
	}

	return nil
}
