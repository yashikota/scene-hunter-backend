package handler

import (
	"net/http"

	"github.com/yashikota/scene-hunter-backend/internal/room"
	"github.com/yashikota/scene-hunter-backend/internal/util"
)

func ChangeGameMaster(w http.ResponseWriter, r *http.Request) {
	user, err := util.ParseAndValidateUser(r, 100)
	if err != nil {
		util.JsonResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get and validate the room ID
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

	// Check if the user exists
	exist, err = room.CheckExistUser(roomID, user.ID)
	if err != nil {
		util.JsonResponse(w, http.StatusInternalServerError, "Failed to check if the user exists")
		return
	}

	if !exist {
		util.JsonResponse(w, http.StatusNotFound, "User does not exist")
		return
	}

	// Change the game master
	err = room.ChangeGameMaster(roomID, user.ID)
	if err != nil {
		util.JsonResponse(w, http.StatusInternalServerError, "Failed to change the game master")
		return
	}

	util.JsonResponse(w, http.StatusOK, "Successfully changed the game master")
}
