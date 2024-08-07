package handler

import (
	"net/http"
	"strconv"

	"github.com/yashikota/scene-hunter-backend/src/util"
)

func GenerateUserIDHandler(w http.ResponseWriter, r *http.Request) {
	ttl, _ := strconv.Atoi(r.URL.Query().Get("ttl"))
	if ttl == 0 { // default user id expiration time is 3 hours
		ttl = 60 * 60 * 3
	}

	userID, err := util.GenerateUserID(ttl)
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusInternalServerError, err)
		return
	}

	util.SuccessJsonResponse(w, http.StatusOK, "user_id", userID)
}

func ExistUserIDHandler(w http.ResponseWriter, r *http.Request) {
	user, err := util.ParseAndValidateUser(r, util.ValidateIDRequired)
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
