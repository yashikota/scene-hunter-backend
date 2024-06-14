package main

import (
	"log"
	"net/http"

	"github.com/rs/cors"
	"github.com/yashikota/scene-hunter-backend/internal/handler"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/generate_user_id", handler.GenerateUserID)
	mux.HandleFunc("GET /api/exist_user_id", handler.ExistUserID)
	mux.HandleFunc("POST /api/create_room", handler.CreateRoomHandler)
	mux.HandleFunc("POST /api/join_room", handler.JoinRoomHandler)
	mux.HandleFunc("DELETE /api/delete_room", handler.DeleteRoomHandler)
	mux.HandleFunc("GET /api/get_room_users", handler.GetRoomUsersHandler)
	mux.HandleFunc("DELETE /api/delete_room_user", handler.DeleteRoomUserHandler)
	mux.HandleFunc("POST /api/upload_photo", handler.UploadPhoto)
	mux.HandleFunc("GET /api/result", handler.ResultHandler)
	mux.HandleFunc("GET /api/ping", handler.PingHandler)

	handler := cors.Default().Handler(mux)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
