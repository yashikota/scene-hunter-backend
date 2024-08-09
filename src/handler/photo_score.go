package handler

import (
	"fmt"
	"net/http"

	"github.com/yashikota/scene-hunter-backend/src/room"
	"github.com/yashikota/scene-hunter-backend/src/util"
)

func PhotoScoreHandler(w http.ResponseWriter, r *http.Request) {
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

	// Check if the user already joined the room
	_, err = room.CheckExistUser(roomID, user.ID)
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusInternalServerError, err)
		return
	}

	gameMasterPhotoUrl, err := room.GetGameMasterPhotoUrl(roomID)
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusInternalServerError, err)
		return
	}

	playerPhotoUrls, err := room.GetPlayerPhotoUrls(roomID, user.ID)
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusInternalServerError, err)
		return
	}

	score, err := util.GetPhotoScore(gameMasterPhotoUrl, playerPhotoUrls[0], playerPhotoUrls[1])
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusInternalServerError, err)
		return
	}

	util.SuccessJsonResponse(w, http.StatusOK, "score", score)
}
