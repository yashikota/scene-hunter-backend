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
	mux.HandleFunc("/api/join_room", handler.JoinRoomHandler)
	mux.HandleFunc("/api/upload_photo", handler.UploadPhoto)
	mux.HandleFunc("/api/result", handler.ResultHandler)
	mux.HandleFunc("/api/delete_room", handler.DeleteRoomHandler)

	handler := cors.Default().Handler(mux)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
