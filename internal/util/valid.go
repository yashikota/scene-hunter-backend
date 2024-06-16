package util

import (
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/yashikota/scene-hunter-backend/model"
)

func ParseAndValidateUser(r *http.Request, checkFlag int) (model.User, error) {
	user := model.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return user, fmt.Errorf("failed to parse the request body")
	}

	const (
		flagID   = 0b100
		flagName = 0b010
		flagLang = 0b001
	)

	if checkFlag&flagID != 0 && user.ID == "" {
		return user, fmt.Errorf("id is required")
	}
	if checkFlag&flagName != 0 && user.Name == "" {
		return user, fmt.Errorf("name is required")
	}
	if checkFlag&flagLang != 0 && user.Lang == "" {
		return user, fmt.Errorf("language is required")
	}

	// Validate UserID
	exist, err := ExistUserID(user.ID)
	if err != nil {
		return user, err
	}
	if !exist {
		return user, fmt.Errorf("invalid user ID")
	}

	// DEBUG: Print the form data
	log.Println("ID:", user.ID, "Name:", user.Name, "Language:", user.Lang)

	return user, nil
}

func ParseAndValidateRoom(r *http.Request) (string, error) {
	roomID := r.URL.Query().Get("room_id")
	if roomID == "" {
		return "", fmt.Errorf("room id is required")
	}

	// DEBUG: Print the form data
	log.Println("Room ID:", roomID)

	return roomID, nil
}

// Load Image decodes the image and returns the image, format, and error
func ValidateFile(fileHeader *multipart.FileHeader) error {
	const maxFileSize = 10 * 1024 * 1024 // 10MB
	var allowedTypes []string = []string{"image/jpeg", "image/png"}

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
