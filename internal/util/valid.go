package util

import (
	"encoding/json"
	"fmt"
	"log"
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
		flagID   = 0b1000
		flagName = 0b0100
		flagLang = 0b0010
		flagStatus = 0b0001
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
	if checkFlag&flagStatus != 0 && user.Status == "" {
		return user, fmt.Errorf("status is required")
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

func ValidateMaxFileSize(r *http.Request) error {
	const maxFileSize = 10 * 1024 * 1024 // 10MB

	r.Body = http.MaxBytesReader(nil, r.Body, maxFileSize)
	if err := r.ParseMultipartForm(maxFileSize); err != nil {
		return fmt.Errorf("file size exceeds the limit")
	}

	return nil
}

func ValidateFileType(r *http.Request) (string, error) {
	file, _, err := r.FormFile("image")
	if err != nil {
		return "", fmt.Errorf("failed to read the form file")
	}
	defer file.Close()

	// Validate the file type
	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		return "", fmt.Errorf("failed to read the form file")
	}

	// Check the file type. jpeg or png
	fileType := http.DetectContentType(buff)
	if fileType != "image/jpeg" && fileType != "image/png" {
		return "", fmt.Errorf("invalid file type")
	}

	return fileType, nil
}
