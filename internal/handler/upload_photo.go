package handler

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/yashikota/scene-hunter-backend/internal/room"
	"github.com/yashikota/scene-hunter-backend/internal/util"
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
	fileType, err := util.ValidateFileType(r)
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusBadRequest, err)
		return
	}

	// multipart.File to bytes.Buffer
	img, err := util.ToBytes(r)
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusBadRequest, err)
		return
	}

	// Read the form file
	fileName, err := util.SaveFile(img, originalPhotoUploadDir, ".jpg")
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusBadRequest, err)
		return
	}

	// DEBUG: Print the form data
	log.Printf("uploadDir: %s, fileName: %s", originalPhotoUploadDir, fileName)

	// Asynchronously process the file and immediately return a response to the client
	go func() {
		img, err := util.LoadImage(img, fileType)
		if err != nil {
			util.ErrorJsonResponse(w, http.StatusInternalServerError, err)
			return
		}

		img = util.Resize(img, 720)

		buf, err := util.ConvertToAVIF(img)
		if err != nil {
			util.ErrorJsonResponse(w, http.StatusInternalServerError, err)
			return
		}

		fileName, err := util.SaveFile(buf, convertedPhotoUploadDir, ".avif")
		if err != nil {
			util.ErrorJsonResponse(w, http.StatusInternalServerError, err)
			return
		}

		imagePath := fmt.Sprintf("%s/%s", convertedPhotoUploadDir, fileName)
		log.Printf("imagePath: %s", imagePath)

		err = room.AddRoomUserPhotoAndScore(roomID, userID, imagePath)
		if err != nil {
			log.Println(err)
			return
		}
	}()

	util.SuccessJsonResponse(w, http.StatusOK, "message", "photo uploaded successfully")
}
