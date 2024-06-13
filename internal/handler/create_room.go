package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/yashikota/scene-hunter-backend/internal/room"
	"github.com/yashikota/scene-hunter-backend/internal/util"
)

func CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	user, err := util.ParseAndValidateUser(r)
	if err != nil {
		util.JsonResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get Room ID
	var roomID string
	for i := 0; i < 10; i++ {
		digits := 6
		roomID, err = util.GenerateRoomID(digits, user.ID)
		if err != nil {
			util.JsonResponse(w, http.StatusInternalServerError, "Failed to generate room ID")
			return
		}

		// Check if the room already exists
		exist := room.CheckExistRoom(roomID)
		if !exist {
			break
		}

		// if 10 times looped, return error
		if i == 9 {
			util.JsonResponse(w, http.StatusInternalServerError, "Failed to generate room ID (10 times looped)")
			return
		}
	}

	// DEBUG: Print the room ID
	log.Println("Room ID:", roomID)

	// Create a room
	err = room.CreateRoom(roomID, user)
	if err != nil {
		util.JsonResponse(w, http.StatusInternalServerError, "Failed to creation room")
		return
	}

	// response
	res := struct {
		RoomID  string `json:"room_id,omitempty"`
		Message string `json:"message"`
	}{
		RoomID:  roomID,
		Message: "Room created successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
