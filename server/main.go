package main

import (
	"log"
	"net/http"

	"github.com/rs/cors"
	"github.com/yashikota/scene-hunter-backend/internal/handler"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/ping", handler.PingHandler)
	mux.HandleFunc("/api/create_room", handler.CreateRoomHandler)
	mux.HandleFunc("/api/upload_photo", handler.UploadPhoto)
	mux.HandleFunc("/api/result", handler.ResultHandler)

	handler := cors.Default().Handler(mux)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
