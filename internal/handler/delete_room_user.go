package handler

import (
	"net/http"

	"github.com/yashikota/scene-hunter-backend/internal/room"
	"github.com/yashikota/scene-hunter-backend/internal/util"
)

func DeleteRoomUserHandler(w http.ResponseWriter, r *http.Request) {
	user, err := util.ParseAndValidateUser(r, 100)
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusBadRequest, err)
		return
	}

	roomID, err := util.ParseAndValidateRoom(r)
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusBadRequest, err)
		return
	}

	// Check if the room exists
	_, err = room.CheckExistRoom(roomID)
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusNotFound, err)
		return
	}

	// Check if the user exists
	_, statusCode, err := room.CheckExistUser(roomID, user.ID)
	if err != nil {
		util.ErrorJsonResponse(w, statusCode, err)
		return
	}

	// Delete the user from the room
	err = room.DeleteRoomUser(roomID, user.ID)
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusInternalServerError, err)
		return
	}

	util.SuccessJsonResponse(w, http.StatusOK, "deleted the user from the room")
}
