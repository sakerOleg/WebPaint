package handlers

import (
	"database/sql"
	"html/template"
)

type Handler struct {
	DB        *sql.DB
	Templates *template.Template
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		DB:        db,
		Templates: template.Must(template.ParseGlob("templates/*.html")),
	}
}

type PageData struct {
	Title      string
	Error      string
	Success    string
	Registered bool
	Username   string
}
