package handlers

import (
	"net/http"
)

func (h *Handler) indexHandler(w http.ResponseWriter, r *http.Request) {
	if !h.isAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	cookie, _ := r.Cookie("session")
	data := PageData{
		Title:    "Главная",
		Username: cookie.Value,
	}

	if redirectCookie, err := r.Cookie("redirect_to_editor"); err == nil && redirectCookie.Value == "true" {
		http.SetCookie(w, &http.Cookie{
			Name:     "redirect_to_editor",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		})
		http.Redirect(w, r, "/editor", http.StatusSeeOther)
		return
	}

	h.Templates.ExecuteTemplate(w, "welcome.html", data)
}

func (h *Handler) editorHandler(w http.ResponseWriter, r *http.Request) {
	if !h.isAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	cookie, _ := r.Cookie("session")
	data := struct {
		PageData
		Username string
	}{
		PageData: PageData{
			Title:      "Пиксельный редактор",
			Registered: true,
		},
		Username: cookie.Value,
	}

	h.Templates.ExecuteTemplate(w, "editor.html", data)
}
