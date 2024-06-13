package handler

import (
	"log"
	"net/http"

	"github.com/yashikota/scene-hunter-backend/internal/room"
	"github.com/yashikota/scene-hunter-backend/internal/util"
)

func JoinRoomHandler(w http.ResponseWriter, r *http.Request) {
	user, err := util.ParseAndValidateUser(r)
	if err != nil {
		util.JsonResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Check if the room exists
	roomID := r.URL.Query().Get("room_id")
	exist := room.CheckExistRoom(roomID)
	if !exist {
		util.JsonResponse(w, http.StatusNotFound, "Room does not exist")
		return
	}

	// Join the room
	err = room.JoinRoom(roomID, user)
	if err != nil {
		util.JsonResponse(w, http.StatusInternalServerError, "Failed to join room")
		return
	}

	// DEBUG: Print the room ID
	log.Println("Room ID:", roomID)

	util.JsonResponse(w, http.StatusOK, "Successfully joined the room")
}
