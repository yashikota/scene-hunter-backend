package util

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/yashikota/scene-hunter-backend/model"
)

func ParseAndValidateUser(r *http.Request, validateFuncs ...func(user model.User) error) (model.User, error) {
	user := model.User{}

	// Parse the request body into the user struct
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return user, fmt.Errorf("failed to parse the request body: %w", err)
	}

	// Apply each validation function
	for _, validate := range validateFuncs {
		if err := validate(user); err != nil {
			return user, err
		}
	}

	// DEBUG: Print the form data
	log.Printf("ID: %s, Name: %s, Language: %s", user.ID, user.Name, user.Lang)

	return user, nil
}

// Individual validation functions
func ValidateIDRequired(user model.User) error {
	if user.ID == "" {
		return fmt.Errorf("id is required")
	}
	exist, err := ExistUserID(user.ID)
	if err != nil {
		return fmt.Errorf("failed to validate user ID: %w", err)
	}
	if !exist {
		return fmt.Errorf("invalid user ID")
	}

	return nil
}

func ValidateNameRequired(user model.User) error {
	if user.Name == "" {
		return fmt.Errorf("name is required")
	}

	return nil
}

func ValidateLangRequired(user model.User) error {
	if user.Lang == "" {
		return fmt.Errorf("language is required")
	}

	return nil
}

func ValidateStatusRequired(user model.User) error {
	if user.Status == "" {
		return fmt.Errorf("status is required")
	}

	return nil
}

func ValidateMaxFileSize(r *http.Request) error {
	const maxFileSize = 5 * 1024 * 1024 // 5MB

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
