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

func ParseAndValidateUser(r *http.Request, checkFlag int) (model.User, error) {
	user := model.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return user, err
	}

	const (
		flagID   = 0b100
		flagName = 0b010
		flagLang = 0b001
	)

	if checkFlag&flagID != 0 && user.ID == "" {
		return user, errors.New("id is required")
	}
	if checkFlag&flagName != 0 && user.Name == "" {
		return user, errors.New("name is required")
	}
	if checkFlag&flagLang != 0 && user.Lang == "" {
		return user, errors.New("lang is required")
	}

	// Validate UserID
	exist, err := ExistUserID(user.ID)
	if err != nil {
		return user, err
	}
	if !exist {
		return user, errors.New("invalid user ID")
	}

	// DEBUG: Print the form data
	log.Println("ID:", user.ID, "Name:", user.Name, "Language:", user.Lang)

	return user, nil
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
