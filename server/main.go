package main

import (
	"log"
	"net/http"

	"github.com/rs/cors"
	"github.com/yashikota/scene-hunter-backend/internal/handler"
)

func main() {
	mux := http.NewServeMux()

	// User ID
	mux.HandleFunc("GET /api/generate_user_id", handler.GenerateUserIDHandler)
	mux.HandleFunc("POST /api/exist_user_id", handler.ExistUserIDHandler)

	// Room
	mux.HandleFunc("POST /api/create_room", handler.CreateRoomHandler)
	mux.HandleFunc("POST /api/join_room", handler.JoinRoomHandler)
	mux.HandleFunc("GET /api/get_room_users", handler.GetRoomUsersHandler)
	mux.HandleFunc("PUT /api/change_game_master", handler.ChangeGameMasterHandler)
	mux.HandleFunc("DELETE /api/delete_room_user", handler.DeleteRoomUserHandler)
	mux.HandleFunc("DELETE /api/delete_room", handler.DeleteRoomHandler)

	// Game
	mux.HandleFunc("POST /api/upload_photo", handler.UploadPhotoHandler)
	mux.HandleFunc("GET /api/result", handler.ResultHandler)

	// Debug
	mux.HandleFunc("GET /api/ping", handler.PingHandler)

	// Swagger UI
	fileServer := http.FileServer(http.Dir("./swagger"))
	mux.Handle("/swagger/", http.StripPrefix("/swagger", fileServer))

	handler := cors.Default().Handler(mux)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
