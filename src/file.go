package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

func saveImage(buf *bytes.Buffer, dirName string, fileExtension string) (string, error) {
	fileName := uuid.New().String()
	fileName += fileExtension

	err := saveFile(buf, dirName, fileName)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func saveFile(buf *bytes.Buffer, dirName, fileName string) error {
	// Create the directory if it does not exist
	err := os.MkdirAll(dirName, 0755)
	if err != nil {
		log.Fatal(err)
	}

	// Save the file
	file, err := os.Create(filepath.Join(dirName, fileName))
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
