package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/yashikota/scene-hunter-backend/internal/util"
	"github.com/yashikota/scene-hunter-backend/model"
)

func CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	user := model.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		util.JsonErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate the parsed data
	if user.ID == "" {
		util.JsonErrorResponse(w, http.StatusBadRequest, "id is required")
		return
	}

	if user.Name == "" {
		util.JsonErrorResponse(w, http.StatusBadRequest, "name is required")
		return
	}

	if user.Lang == "" {
		util.JsonErrorResponse(w, http.StatusBadRequest, "lang is required")
		return
	}

	// DEBUG: Print the form data
	log.Println("ID:", user.ID, "Name:", user.Name, "Language:", user.Lang)

	// Get Room ID
	digits := 6
	roomID, err := util.GenerateRoomID(digits, user.ID)
	if err != nil {
		util.JsonErrorResponse(w, http.StatusInternalServerError, "Failed to creation room")
		return
	}

	// DEBUG: Print the room ID
	log.Println("Room ID:", roomID)

	// response
	res := model.Room{
		RoomID: roomID,
		Message: "Room created successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
