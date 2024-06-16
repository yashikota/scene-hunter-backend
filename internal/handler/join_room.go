package handler

import (
	"log"
	"net/http"

	"github.com/yashikota/scene-hunter-backend/internal/room"
	"github.com/yashikota/scene-hunter-backend/internal/util"
)

func JoinRoomHandler(w http.ResponseWriter, r *http.Request) {
	user, err := util.ParseAndValidateUser(r, 111) // Validate ID, Name, and Language
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusBadRequest, err)
		return
	}

	// Check if the room exists
	roomID, err := util.ParseAndValidateRoom(r)
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusBadRequest, err)
		return
	}

	_, err = room.CheckExistRoom(roomID)
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusNotFound, err)
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

	util.SuccessJsonResponse(w, http.StatusOK, "successfully joined the room")
}
