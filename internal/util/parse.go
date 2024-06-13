package util

import (
	"encoding/json"
	"errors"
	"log"
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
