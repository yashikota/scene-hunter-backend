package handler

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/yashikota/scene-hunter-backend/src/room"
	"github.com/yashikota/scene-hunter-backend/src/util"
)

func UploadPhotoHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(1024)
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusBadRequest, err)
		return
	}

	// Parse and validate the user
	userID := r.FormValue("user_id")
	if userID == "" {
		util.ErrorJsonResponse(w, http.StatusBadRequest, fmt.Errorf("user_id is required"))
		return
	}
	// Validate UserID
	exist, err := util.ExistUserID(userID)
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusInternalServerError, err)
		return
	}
	if !exist {
		util.ErrorJsonResponse(w, http.StatusBadRequest, fmt.Errorf("invalid user ID"))
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

	// Create directory
	originalPhotoUploadDir := filepath.Join("uploads", roomID, "original")
	err = util.MakeDir(originalPhotoUploadDir)
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusInternalServerError, err)
		return
	}

	convertedPhotoUploadDir := filepath.Join("uploads", roomID, "converted")
	err = util.MakeDir(convertedPhotoUploadDir)
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusInternalServerError, err)
		return
	}

	// Validate Max File Size
	err = util.ValidateMaxFileSize(r)
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusBadRequest, err)
		return
	}

	// Validate File Type
	_, err = util.ValidateFileType(r)
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusBadRequest, err)
		return
	}

	// multipart.File to bytes.Buffer
	img, err := util.ReadFormFile(r)
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusBadRequest, err)
		return
	}

	// Save the original photo
	originalFileName, err := util.SaveFile(img, originalPhotoUploadDir, ".jpg")
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusBadRequest, err)
		return
	}
	originalFilePath := fmt.Sprintf("%s/%s", originalPhotoUploadDir, originalFileName)

	// Resize the photo
	resizedImg := util.Resize(img, 720)
	convertedFileName, err := util.SaveFile(resizedImg, convertedPhotoUploadDir, ".jpg")
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusInternalServerError, err)
		return
	}
	convertedFilePath := fmt.Sprintf("%s/%s", convertedPhotoUploadDir, convertedFileName)

	uploadLog := map[string]string{
		"original":  originalFilePath,
		"converted": convertedFilePath,
	}

	err = room.AddRoomUserPhotoAndScore(roomID, userID, convertedFilePath)
	if err != nil {
		log.Println(err)
		return
	}

	util.SuccessJsonResponse(w, http.StatusOK, "image_path", uploadLog)
}
