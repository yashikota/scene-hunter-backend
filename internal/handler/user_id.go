package handler

import (
	"net/http"

	"github.com/yashikota/scene-hunter-backend/internal/util"
)

func GenerateUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := util.GenerateUserID()
	if err != nil {
		util.JsonResponse(w, http.StatusInternalServerError, "Failed to generate user ID")
		return
	}

	util.JsonResponse(w, http.StatusOK, userID)
}

func ExistUserID(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		util.JsonResponse(w, http.StatusBadRequest, "user_id is required")
		return
	}

	exist, err := util.ExistUserID(userID)
	if err != nil {
		util.JsonResponse(w, http.StatusInternalServerError, "Failed to check user ID")
		return
	}

	if !exist {
		util.JsonResponse(w, http.StatusNotFound, "User ID does not exist")
		return
	}

	util.JsonResponse(w, http.StatusOK, "User ID exists")
}
