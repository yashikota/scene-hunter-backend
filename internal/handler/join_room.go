package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/yashikota/scene-hunter-backend/internal/room"
	"github.com/yashikota/scene-hunter-backend/internal/util"
	"github.com/yashikota/scene-hunter-backend/model"
)

func JoinRoomHandler(w http.ResponseWriter, r *http.Request) {
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
