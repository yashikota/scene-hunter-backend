package handler

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/yashikota/scene-hunter-backend/internal/room"
	"github.com/yashikota/scene-hunter-backend/internal/util"
)

func UploadPhotoHandler(w http.ResponseWriter, r *http.Request) {
	const uploadDir = "uploads"
	const maxConcurrentWorkers = 5

	user, err := util.ParseAndValidateUser(r, 100)
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusBadRequest, err)
		return
	}
	log.Printf("User ID: %s", user.ID)

	// Check if the room exists
	roomID := r.URL.Query().Get("room_id")
	_, err = room.CheckExistRoom(roomID)
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusNotFound, err)
		return
	}

	// Create RoomID directory
	err = util.MakeDir(uploadDir)
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusInternalServerError, err)
		return
	}
	roomDirPath := filepath.Join(uploadDir, roomID)

	fileHeaders, err := util.GetFileHeaders(r)
	if err != nil {
		util.ErrorJsonResponse(w, http.StatusBadRequest, err)
		return
	}

	jobs, results := make(chan *multipart.FileHeader, len(fileHeaders)), make(chan string, len(fileHeaders))
	var wg sync.WaitGroup

	for i := 0; i < maxConcurrentWorkers; i++ {
		wg.Add(1)
		go worker(jobs, results, roomDirPath, &wg)
	}

	for _, fileHeader := range fileHeaders {
		if err := util.ValidateFile(fileHeader); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		jobs <- fileHeader
	}

	close(jobs)
	wg.Wait()
	close(results)

	writeResults(w, results)
}

func writeResults(w http.ResponseWriter, results <-chan string) {
	for url := range results {
		w.Write([]byte(url + "\n"))
	}
}

func worker(jobs <-chan *multipart.FileHeader, results chan<- string, roomDirPath string, wg *sync.WaitGroup) {
	defer wg.Done()
	for fileHeader := range jobs {
		processFile(fileHeader, results, roomDirPath)
	}
}

func processFile(fileHeader *multipart.FileHeader, results chan<- string, roomDirPath string) {
	file, err := fileHeader.Open()
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return
	}
	defer file.Close()

	img, format, err := util.LoadImage(file)
	if err != nil {
		log.Printf("Error decoding image: %v", err)
		return
	}

	log.Printf("file name: %s, format: %s", fileHeader.Filename, format)

	img = util.Resize(img, 720)

	buf, err := util.ConvertToAVIF(img)
	if err != nil {
		log.Printf("Error converting to AVIF: %v", err)
		return
	}

	fileExtension := ".avif"
	fileName, err := util.SaveImage(buf, roomDirPath, fileExtension)
	if err != nil {
		log.Printf("Error saving file: %v", err)
		return
	}

	url := fmt.Sprintf("http://%s/%s/%s", "localhost:8080", roomDirPath, fileName)
	results <- url
}
