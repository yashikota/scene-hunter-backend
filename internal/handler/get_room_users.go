package handler

import (
	"net/http"

	"github.com/yashikota/scene-hunter-backend/internal/room"
	"github.com/yashikota/scene-hunter-backend/internal/util"
)

func GetRoomUsersHandler(w http.ResponseWriter, r *http.Request) {
	roomID, err := util.ParseAndValidateRoom(r)
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusBadRequest, err)
		return
	}

	// Get users in the room
	users, err := room.GetRoomUsers(roomID)
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusInternalServerError, err)
		return
	}

	// Response
	util.SuccessJsonResponse(w, http.StatusOK, "users", users)
}
