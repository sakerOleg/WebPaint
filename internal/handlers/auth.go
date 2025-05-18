package handlers

import (
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) loginHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:      "Вход",
		Registered: r.URL.Query().Get("registered") == "true",
	}

	if r.Method == "GET" {
		h.Templates.ExecuteTemplate(w, "login.html", data)
		return
	}

	if err := r.ParseForm(); err != nil {
		data.Error = "Неверные данные формы"
		h.Templates.ExecuteTemplate(w, "login.html", data)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	var user struct {
		ID       int
		Username string
		Password string
	}

	err := h.DB.QueryRow("SELECT id, username, password FROM users WHERE username = ?", username).
		Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		data.Error = "Неверный логин или пароль"
		h.Templates.ExecuteTemplate(w, "login.html", data)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		data.Error = "Неверный логин или пароль"
		h.Templates.ExecuteTemplate(w, "login.html", data)
		return
	}

	h.createSession(w, user.Username)

	http.SetCookie(w, &http.Cookie{
		Name:     "redirect_to_editor",
		Value:    "true",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(3 * time.Second),
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) registerHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title: "Регистрация",
	}

	if r.Method == "GET" {
		h.Templates.ExecuteTemplate(w, "register.html", data)
		return
	}

	if err := r.ParseForm(); err != nil {
		data.Error = "Неверные данные формы"
		h.Templates.ExecuteTemplate(w, "register.html", data)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	confirm := r.FormValue("confirm_password")

	if password != confirm {
		data.Error = "Пароли не совпадают"
		h.Templates.ExecuteTemplate(w, "register.html", data)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		data.Error = "Ошибка сервера"
		h.Templates.ExecuteTemplate(w, "register.html", data)
		return
	}

	_, err = h.DB.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, string(hashedPassword))
	if err != nil {
		data.Error = "Пользователь с таким именем уже существует"
		h.Templates.ExecuteTemplate(w, "register.html", data)
		return
	}

	http.Redirect(w, r, "/login?registered=true", http.StatusSeeOther)
}

func (h *Handler) logoutHandler(w http.ResponseWriter, r *http.Request) {
	h.clearSession(w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *Handler) createSession(w http.ResponseWriter, username string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    username,
		Path:     "/",
		HttpOnly: true,
	})
}

func (h *Handler) clearSession(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
}

func (h *Handler) isAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie("session")
	return err == nil && cookie.Value != ""
}
