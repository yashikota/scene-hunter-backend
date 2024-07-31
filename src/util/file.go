package util

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/oklog/ulid/v2"
)

func MakeDir(dirName string) error {
	// If the directory does not exist, create it
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err := os.MkdirAll(dirName, 0755)
		if err != nil {
			return fmt.Errorf("failed to create the directory: %w", err)
		}
	}

	return nil
}

func DeleteDir(dirName string) error {
	// Delete the directory
	err := os.RemoveAll(dirName)
	if err != nil {
		return err
	}

	return nil
}

func SaveFile(data []byte, path string, extension string) (string, error) {
	// Generate a unique file name
	ulid, err := ulid.New(ulid.Now(), nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate the unique file name: %w", err)
	}
	fileName := ulid.String()
	fileName += extension

	// Save the file
	file, err := os.Create(filepath.Join(path, fileName))
	if err != nil {
		return "", fmt.Errorf("failed to save the file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("failed to save the file: %w", err)
	}

	return fileName, nil
}

func ReadFormFile(r *http.Request) ([]byte, error) {
	file, _, err := r.FormFile("image")
	if err != nil {
		return nil, fmt.Errorf("failed to read the form file")
	}
	defer file.Close()

	// Read the form file
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, file)
	if err != nil {
		return nil, fmt.Errorf("failed to read the form file")
	}

	return buf.Bytes(), nil
}
