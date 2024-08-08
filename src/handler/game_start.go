package handler

import (
	"fmt"
	"net/http"

	"github.com/yashikota/scene-hunter-backend/src/room"
	"github.com/yashikota/scene-hunter-backend/src/util"
	"github.com/yashikota/scene-hunter-backend/model"
)

func GameStartHandler(w http.ResponseWriter, r* http.Request) {
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

	result, err := room.CheckExistRoom(roomID)
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusInternalServerError, err)
		return
	}
	if !result {
		util.ErrorJsonResponse(w, http.StatusNotFound, err)
		return
	}

	// Check user is the game master
	isGameMaster, status, err := room.CheckGameMaster(roomID, user.ID)
	if !isGameMaster {
		util.ErrorJsonResponse(w, status, err)
		return
	}

	// Start the game
	err = room.UpdateGameStatus(roomID, model.GameMasterPhoto)
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusInternalServerError, err)
		return
	}

	util.SuccessJsonResponse(w, http.StatusOK, "message", "successfully started the game")
}
