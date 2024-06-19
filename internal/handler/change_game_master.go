package handler

import (
	"fmt"
	"net/http"

	"github.com/yashikota/scene-hunter-backend/internal/room"
	"github.com/yashikota/scene-hunter-backend/internal/util"
)

func ChangeGameMasterHandler(w http.ResponseWriter, r *http.Request) {
	user, err := util.ParseAndValidateUser(r, 1000) // Validate ID
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusBadRequest, err)
		return
	}

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


	// Check if the user exists
	_, statusCode, err := room.CheckExistUser(roomID, user.ID)
	if err != nil {
		util.ErrorJsonResponse(w, statusCode, err)
		return
	}

	// Change the game master
	err = room.ChangeGameMaster(roomID, user.ID)
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusInternalServerError, err)
		return
	}

	util.SuccessJsonResponse(w, http.StatusOK, "game master id", user.ID)
}
