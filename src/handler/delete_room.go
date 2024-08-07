package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/yashikota/scene-hunter-backend/src/room"
	"github.com/yashikota/scene-hunter-backend/src/util"
)

func DeleteRoomHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the room exists
	roomID := r.URL.Query().Get("room_id")
	if roomID == "" {
		util.ErrorJsonResponse(w, http.StatusBadRequest, fmt.Errorf("room_id is required"))
		return
	}

	result, err := room.CheckExistRoom(roomID)
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusInternalServerError, err)
		return
	}
	if !result {
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
