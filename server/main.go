package main

import (
	"log"
	"net/http"

	"github.com/rs/cors"
	"github.com/yashikota/scene-hunter-backend/internal/handler"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/ping", handler.PingHandler)
	mux.HandleFunc("POST /api/create_room", handler.CreateRoomHandler)
	mux.HandleFunc("POST /api/join_room", handler.JoinRoomHandler)
	mux.HandleFunc("GET /api/get_room_users", handler.GetRoomUsersHandler)
	mux.HandleFunc("POST /api/upload_photo", handler.UploadPhoto)
	mux.HandleFunc("GET /api/result", handler.ResultHandler)
	mux.HandleFunc("DELETE /api/delete_room", handler.DeleteRoomHandler)

	handler := cors.Default().Handler(mux)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
