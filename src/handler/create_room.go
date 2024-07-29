package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/yashikota/scene-hunter-backend/src/room"
	"github.com/yashikota/scene-hunter-backend/src/util"
)

func CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	user, err := util.ParseAndValidateUser(r, 1110) // Validate ID, Name, and Language
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusBadRequest, err)
		return
	}

	// Get Room ID
	var roomID string
	for i := 0; i < 10; i++ {
		digits := 6
		roomID, err = util.GenerateRoomID(digits, user.ID)
		if err != nil {
			util.ErrorJsonResponse(w, http.StatusInternalServerError, err)
			return
		}

		// Check if the room already exists
		exist, _ := room.CheckExistRoom(roomID)
		if !exist { // if not exists
			break
		}

		// if 10 times looped, return error
		if i == 9 {
			util.ErrorJsonResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to generate room id"))
			return
		}
	}

	// DEBUG: Print the room ID
	log.Println("Room ID:", roomID)

	// Create a room
	err = room.CreateRoom(roomID, user)
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusInternalServerError, err)
		return
	}

	util.SuccessJsonResponse(w, http.StatusCreated, "room_id", roomID)
}
