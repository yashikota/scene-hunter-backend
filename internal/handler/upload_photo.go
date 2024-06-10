package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/yashikota/scene-hunter-backend/internal/util"
)

func UploadPhoto(w http.ResponseWriter, r *http.Request) {
	const uploadDir = "uploads"

	// Create a multipart reader
	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Search image file
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if part.FileName() != "" {
			// Decode the image
			img, format, err := util.LoadImage(part)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// DEBUG: Print the file name and format
			log.Println("file name:", part.FileName(), "format:", format)

			// Resize the image to HD
			img = util.Resize(img, 720)

			// Convert to AVIF
			buf, err := util.ConvertToAVIF(img)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Save the AVIF image
			fileExtension := ".avif"
			filename, err := util.SaveImage(buf, uploadDir, fileExtension)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Return the URL of the uploaded image
			url := fmt.Sprintf("http://%s/%s/%s", r.Host, uploadDir, filename)
			w.Write([]byte(url + "\n"))
			return
		}
	}

	http.Error(w, "Image not found", http.StatusBadRequest)
}
