package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/yashikota/scene-hunter-backend/model"
)

func ParseAndValidateUser(r *http.Request) (model.User, error) {
	user := model.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return user, err
	}

	if user.ID == "" {
		return user, errors.New("id is required")
	}
	if user.Name == "" {
		return user, errors.New("name is required")
	}
	if user.Lang == "" {
		return user, errors.New("lang is required")
	}

	// DEBUG: Print the form data
	log.Println("ID:", user.ID, "Name:", user.Name, "Language:", user.Lang)

	return user, nil
}

// Upload Photo validates the request and returns the room ID, user ID, and error
func ParseAndValidateRequest(r *http.Request) (string, string, error) {
	roomID := r.URL.Query().Get("room_id")
	if roomID == "" {
		return "", "", fmt.Errorf("room ID is required")
	}

	userID := r.FormValue("user_id")
	if userID == "" {
		return "", "", fmt.Errorf("user ID is required")
	}

	return roomID, userID, nil
}

// Load Image decodes the image and returns the image, format, and error
func ValidateFile(fileHeader *multipart.FileHeader) error {
	const maxFileSize = 10 * 1024 * 1024 // 10MB
	var allowedTypes []string = []string{"image/jpeg", "image/jpg", "image/png"}

	if fileHeader.Size > maxFileSize {
		return fmt.Errorf("file size exceeds the limit")
	}

	detectedType := http.DetectContentType([]byte(fileHeader.Header.Get("Content-Type")))
	log.Println("Detected file type:", detectedType)
	for _, allowedType := range allowedTypes {
		if detectedType == allowedType {
			return nil
		}
	}

	return fmt.Errorf("invalid file type")
}
