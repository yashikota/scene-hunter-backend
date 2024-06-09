package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("POST /upload", uploadHandler)

	http.HandleFunc("GET /uploads/", uploadsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

const uploadDir = "uploads"

func uploadHandler(w http.ResponseWriter, r *http.Request) {
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
			img, format, err := loadImage(part)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// DEBUG: Print the file name and format
			log.Println("file name:", part.FileName(), "format:", format)

			// Resize the image to HD
			img = resize(img, 720)

			// Convert to AVIF
			buf, err := convertToAVIF(img)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Save the AVIF image
			fileExtension := ".avif"
			filename, err := saveImage(buf, uploadDir, fileExtension)
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

func uploadsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/uploads/" {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}

	// Serve the files under /uploads/ directory
	http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadDir))).ServeHTTP(w, r)
}
