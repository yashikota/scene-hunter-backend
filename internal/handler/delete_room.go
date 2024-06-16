package handler

import (
	"log"
	"net/http"

	"github.com/yashikota/scene-hunter-backend/internal/room"
	"github.com/yashikota/scene-hunter-backend/internal/util"
)

func DeleteRoomHandler(w http.ResponseWriter, r *http.Request) {
	roomID, err := util.ParseAndValidateRoom(r)
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusBadRequest, err)
		return
	}

	// Check if the room exists
	_, err = room.CheckExistRoom(roomID)
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusNotFound, err)
		return
	}

	// Delete the room
	err = room.DeleteRoom(roomID)
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusInternalServerError, err)
		return
	}

	// DEBUG: Print the room ID
	log.Println("Room ID:", roomID)

	util.SuccessJsonResponse(w, http.StatusOK, "message", "successfully deleted the room")
}
