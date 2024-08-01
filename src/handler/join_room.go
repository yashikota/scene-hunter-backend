package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/yashikota/scene-hunter-backend/src/room"
	"github.com/yashikota/scene-hunter-backend/src/util"
)

func JoinRoomHandler(w http.ResponseWriter, r *http.Request) {
	user, err := util.ParseAndValidateUser(r, util.ValidateIDRequired, util.ValidateNameRequired, util.ValidateLangRequired)
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusBadRequest, err)
		return
	}

	// Check if the room exists
	roomID := r.URL.Query().Get("room_id")
	if roomID == "" {
		util.ErrorJsonResponse(w, http.StatusBadRequest, fmt.Errorf("room_id is required"))
		return
	}

	// Check if the room exists
	result, err := room.CheckExistRoom(roomID)
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusInternalServerError, err)
		return
	}
	if !result {
		util.ErrorJsonResponse(w, http.StatusNotFound, err)
		return
	}

	// Check if the user already joined the room
	exist, status, err := room.CheckExistUser(roomID, user.ID)
	if err != nil {
		util.ErrorJsonResponse(w, status, err)
		return
	}
	if exist {
		util.ErrorJsonResponse(w, http.StatusConflict, fmt.Errorf("user already joined the room"))
		return
	}

	// Join the room
	err = room.JoinRoom(roomID, user)
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusInternalServerError, err)
		return
	}

	// DEBUG: Print the room ID
	log.Println("Room ID:", roomID)

	util.SuccessJsonResponse(w, http.StatusOK, "message", "successfully joined the room")
}
