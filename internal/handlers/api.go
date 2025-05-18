package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func (h *Handler) saveDrawingHandler(w http.ResponseWriter, r *http.Request) {
	if !h.isAuthenticated(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	cookie, _ := r.Cookie("session")
	username := cookie.Value

	var userID int
	err := h.DB.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	var data struct {
		Title string `json:"title"`
		Image string `json:"image"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	_, err = h.DB.Exec("INSERT INTO drawings (user_id, title, data) VALUES (?, ?, ?)",
		userID, data.Title, data.Image)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) getDrawingsHandler(w http.ResponseWriter, r *http.Request) {
	if !h.isAuthenticated(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	cookie, _ := r.Cookie("session")
	username := cookie.Value

	var userID int
	err := h.DB.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	rows, err := h.DB.Query(`
        SELECT id, title, data as image, created_at 
        FROM drawings 
        WHERE user_id = ? 
        ORDER BY created_at DESC`, userID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Drawing struct {
		ID        int       `json:"id"`
		Title     string    `json:"title"`
		Image     string    `json:"image"`
		CreatedAt time.Time `json:"created_at"`
	}

	var drawings []Drawing
	for rows.Next() {
		var d Drawing
		rows.Scan(&d.ID, &d.Title, &d.Image, &d.CreatedAt)
		drawings = append(drawings, d)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(drawings)
}

func (h *Handler) deleteDrawingHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Проверка аутентификации
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Error(w, "Session cookie required", http.StatusUnauthorized)
		return
	}

	// 2. Получение user_id
	var userID int
	err = h.DB.QueryRow("SELECT id FROM users WHERE username = ?", cookie.Value).Scan(&userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	// 3. Извлечение ID рисунка из URL
	drawingID := strings.TrimPrefix(r.URL.Path, "/api/delete-drawing/")
	if drawingID == "" {
		http.Error(w, "Drawing ID is required", http.StatusBadRequest)
		return
	}

	// 4. Проверка что ID - число
	var drawingIDInt int
	if _, err := fmt.Sscanf(drawingID, "%d", &drawingIDInt); err != nil {
		http.Error(w, "Invalid drawing ID format", http.StatusBadRequest)
		return
	}

	// 5. Выполнение удаления с проверкой владельца
	result, err := h.DB.Exec(`
        DELETE FROM drawings 
        WHERE id = ? AND user_id = ?`,
		drawingIDInt, userID)

	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 6. Проверка что запись действительно удалена
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to check deletion", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Drawing not found or not owned by user", http.StatusNotFound)
		return
	}

	// 7. Успешный ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":     "success",
		"deleted_id": drawingIDInt,
	})
}
