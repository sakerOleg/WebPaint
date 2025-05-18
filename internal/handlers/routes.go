package handlers

import (
	"net/http"
)

func SetupRoutes(h *Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", h.indexHandler)
	mux.HandleFunc("/login", h.loginHandler)
	mux.HandleFunc("/register", h.registerHandler)
	mux.HandleFunc("/logout", h.logoutHandler)
	mux.HandleFunc("/editor", h.editorHandler)
	mux.HandleFunc("/api/save-drawing", h.saveDrawingHandler)
	mux.HandleFunc("/api/get-drawings", h.getDrawingsHandler)
	mux.HandleFunc("/api/delete-drawing/", h.deleteDrawingHandler)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	return mux
}
