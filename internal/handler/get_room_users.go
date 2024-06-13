package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/yashikota/scene-hunter-backend/internal/room"
	"github.com/yashikota/scene-hunter-backend/internal/util"
)

func GetRoomUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Get the room ID
	roomID := r.URL.Query().Get("room_id")

	// Check if the room exists
	exist := room.CheckExistRoom(roomID)
	if !exist {
		util.JsonResponse(w, http.StatusNotFound, "Room does not exist")
		return
	}

	// Get users in the room
	users, err := room.GetRoomUsers(roomID)
	if err != nil {
		util.JsonResponse(w, http.StatusInternalServerError, "Failed to get users in the room")
		return
	}

	// DEBUG: Print the users in the room
	log.Println("Users in the room:", *users)

	// Response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
