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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the parsed data
	if user.ID == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	if user.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	if user.Lang == "" {
		http.Error(w, "Language is required", http.StatusBadRequest)
		return
	}

	// DEBUG: Print the form data
	log.Println("ID:", user.ID, "Name:", user.Name, "Language:", user.Lang)

	// Get Room ID
	digits := 6
	roomID := util.GenerateRoomID(digits, user.ID)

	// DEBUG: Print the room ID
	log.Println("Room ID:", roomID)
}
