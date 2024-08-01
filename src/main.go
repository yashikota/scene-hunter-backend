package main

import (
	"net/http"
	"strings"

	"github.com/yashikota/scene-hunter-backend/src/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// User ID
	r.HandleFunc("GET /api/generate_user_id", handler.GenerateUserIDHandler)
	r.HandleFunc("POST /api/exist_user_id", handler.ExistUserIDHandler)

	// Room
	r.HandleFunc("POST /api/create_room", handler.CreateRoomHandler)
	r.HandleFunc("POST /api/join_room", handler.JoinRoomHandler)
	r.HandleFunc("GET /api/get_room_users", handler.GetRoomUsersHandler)
	r.HandleFunc("PATCH /api/change_game_master", handler.ChangeGameMasterHandler)
	r.HandleFunc("DELETE /api/delete_room_user", handler.DeleteRoomUserHandler)
	r.HandleFunc("DELETE /api/delete_room", handler.DeleteRoomHandler)
	r.HandleFunc("PUT /api/update_rounds", handler.UpdateRoundsHandler)
	r.HandleFunc("GET /api/get_game_status", handler.GetGameStatusHandler)

	// Game
	r.HandleFunc("POST /api/upload_photo", handler.UploadPhotoHandler)

	// Photo Preview
	photoServer := http.FileServer(http.Dir("./uploads"))
	r.Handle("/uploads/*", http.StripPrefix("/uploads/", photoServer))

	// Debug
	r.HandleFunc("GET /api/ping", handler.PingHandler)

	// Swagger UI
	FileServer(r, "/swagger", http.Dir("./swagger"))

	// Start Server
	http.ListenAndServe(":8080", r)
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusPermanentRedirect).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	})
}
