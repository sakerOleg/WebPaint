package main

import (
	"log"
	"net/http"

	"webpaint/internal/config"
	"webpaint/internal/database"
	"webpaint/internal/handlers"
)

func main() {
	cfg := config.Load()
	db := database.InitDB(cfg)
	defer db.DB.Close()

	handler := handlers.NewHandler(db.DB)
	router := handlers.SetupRoutes(handler)

	log.Printf("Сервер запущен на http://localhost:%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
