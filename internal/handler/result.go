package handler

import (
	"net/http"
)

func ResultHandler(w http.ResponseWriter, r *http.Request) {
	const uploadDir = "uploads"

	if r.URL.Path == "/uploads/" {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}

	// Serve the files under /uploads/ directory
	http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadDir))).ServeHTTP(w, r)
}
