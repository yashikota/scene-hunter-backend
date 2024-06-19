package handler

import (
	"net/http"
	"strconv"

	"github.com/yashikota/scene-hunter-backend/internal/util"
)

func GenerateUserIDHandler(w http.ResponseWriter, r *http.Request) {
	ttl, _ := strconv.Atoi(r.URL.Query().Get("ttl"))
	if ttl == 0 { // if ttl is not set, set 1 day
		ttl = 60 * 60 * 24 // 1 day
	}

	userID, err := util.GenerateUserID(ttl)
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusInternalServerError, err)
		return
	}

	util.SuccessJsonResponse(w, http.StatusOK, "user id", userID)
}

func ExistUserIDHandler(w http.ResponseWriter, r *http.Request) {
	user, err := util.ParseAndValidateUser(r, 1000) // Validate ID
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusBadRequest, err)
		return
	}

	_, err = util.ExistUserID(user.ID)
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusInternalServerError, err)
		return
	}

	util.SuccessJsonResponse(w, http.StatusOK, "message", "user id exists")
}
