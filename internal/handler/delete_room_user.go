package handler

import (
	"net/http"

	"github.com/yashikota/scene-hunter-backend/internal/room"
	"github.com/yashikota/scene-hunter-backend/internal/util"
)

func DeleteRoomUserHandler(w http.ResponseWriter, r *http.Request) {
	user, err := util.ParseAndValidateUser(r)
	if err != nil {
		util.JsonResponse(w, http.StatusBadRequest, err.Error())
		return
	}

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

	// Delete the user from the room
	err = room.DeleteRoomUser(roomID, user.ID)
	if err != nil {
		util.JsonResponse(w, http.StatusInternalServerError, "Failed to delete user from the room")
		return
	}

	util.JsonResponse(w, http.StatusOK, "Successfully deleted the user from the room")
}
