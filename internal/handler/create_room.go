package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/yashikota/scene-hunter-backend/internal/room"
	"github.com/yashikota/scene-hunter-backend/internal/util"
	"github.com/yashikota/scene-hunter-backend/model"
)

func CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	user := model.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		util.JsonResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate the parsed data
	if user.ID == "" {
		util.JsonResponse(w, http.StatusBadRequest, "id is required")
		return
	}

	if user.Name == "" {
		util.JsonResponse(w, http.StatusBadRequest, "name is required")
		return
	}

	if user.Lang == "" {
		util.JsonResponse(w, http.StatusBadRequest, "lang is required")
		return
	}

	// DEBUG: Print the form data
	log.Println("ID:", user.ID, "Name:", user.Name, "Language:", user.Lang)

	// Get Room ID
	var roomID string
	for {
		digits := 6
		roomID, err = util.GenerateRoomID(digits, user.ID)
		if err != nil {
			util.JsonResponse(w, http.StatusInternalServerError, "Failed to creation room")
			return
		}

		// Check if the room already exists
		exist := room.CheckExistRoom(roomID)
		if !exist {
			break
		}
	}

	// DEBUG: Print the room ID
	log.Println("Room ID:", roomID)

	// Create a room
	err = room.CreateRoom(roomID)
	if err != nil {
		util.JsonResponse(w, http.StatusInternalServerError, "Failed to creation room")
		return
	}

	// Join the room
	err = room.JoinRoom(roomID, user)
	if err != nil {
		util.JsonResponse(w, http.StatusInternalServerError, "Failed to join room")
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
