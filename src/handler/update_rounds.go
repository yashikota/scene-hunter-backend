package handler

import (
	"fmt"
	"net/http"

	"github.com/yashikota/scene-hunter-backend/src/room"
	"github.com/yashikota/scene-hunter-backend/src/util"
)

func UpdateRoundsHandler(w http.ResponseWriter, r *http.Request) {
	user, err := util.ParseAndValidateUser(r, util.ValidateIDRequired)
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

	// Check if the room exists
	result, err := room.CheckExistRoom(roomID)
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusInternalServerError, err)
		return
	}
	if !result {
		util.ErrorJsonResponse(w, http.StatusNotFound, err)
		return
	}

	// Check if the user is the game master
	_, status, err := room.CheckGameMaster(roomID, user.ID)
	if err != nil {
		util.ErrorJsonResponse(w, status, err)
		return
	}

	// Get the new rounds
	newRounds := r.URL.Query().Get("rounds")
	if newRounds == "" {
		util.ErrorJsonResponse(w, http.StatusBadRequest, fmt.Errorf("rounds is required"))
		return
	}

	// Update the rounds
	err = room.UpdateRounds(roomID, newRounds)
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusInternalServerError, err)
		return
	}

	util.SuccessJsonResponse(w, http.StatusOK, "new_rounds", newRounds)
}
