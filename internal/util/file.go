package util

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

func MakeDir(dirName string) error {
	// If the directory does not exist, create it
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err := os.MkdirAll(dirName, 0755)
		if err != nil {
			return err
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

func SaveImage(buf *bytes.Buffer, dirName string, fileExtension string) (string, error) {
	fileName := uuid.New().String()
	fileName += fileExtension

	path := filepath.Join(dirName, fileName)

	err := saveFile(buf, path)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func saveFile(buf *bytes.Buffer, path string) error {
	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = io.Copy(file, buf)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func GetFileHeaders(r *http.Request) ([]*multipart.FileHeader, error) {
	const maxFileSize = 10 * 1024 * 1024 // 10MB

	err := r.ParseMultipartForm(maxFileSize)
	if err != nil {
		return nil, err
	}

	return r.MultipartForm.File["image"], nil
}
