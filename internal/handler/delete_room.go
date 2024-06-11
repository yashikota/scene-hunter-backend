package handler

import (
	"log"
	"net/http"

	"github.com/yashikota/scene-hunter-backend/internal/room"
	"github.com/yashikota/scene-hunter-backend/internal/util"
)

func DeleteRoomHandler(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("room_id")
	if roomID == "" {
		util.JsonResponse(w, http.StatusBadRequest, "room_id is required")
		return
	}

	// Check if the room exists
	exist := room.CheckExistRoom(roomID)
	if !exist {
		util.JsonResponse(w, http.StatusNotFound, "Room does not exist")
		return
	}

	// Delete the room
	err := room.DeleteRoom(roomID)
	if err != nil {
		util.JsonResponse(w, http.StatusInternalServerError, "Failed to delete room")
		return
	}

	// DEBUG: Print the room ID
	log.Println("Room ID:", roomID)

	util.JsonResponse(w, http.StatusOK, "Successfully deleted the room")
}
