package handler

import (
	"fmt"
	"net/http"

	"github.com/yashikota/scene-hunter-backend/src/util"
)

func PhotoScoreHandler(w http.ResponseWriter, r *http.Request) {
	gameMasterUrl := r.URL.Query().Get("gamemaster_url")
	player1Url := r.URL.Query().Get("player1_url")
	player2Url := r.URL.Query().Get("player2_url")

	if gameMasterUrl == "" || player1Url == "" || player2Url == "" {
		util.ErrorJsonResponse(w, http.StatusBadRequest, fmt.Errorf("invalid request payload"))
		return
	}

	score, err := util.GetPhotoScore(gameMasterUrl, player1Url, player2Url)
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusInternalServerError, err)
		return
	}

	util.SuccessJsonResponse(w, http.StatusOK, "score", score)
}
