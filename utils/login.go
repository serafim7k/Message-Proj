package utils

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

func LoginHandler(db *sql.DB, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			username := strings.ToLower(r.FormValue("username"))
			password := r.FormValue("password")
			var userID int
			err := db.QueryRow("SELECT id FROM users WHERE LOWER(username) = ? AND password = ?", username, password).Scan(&userID)
			if err != nil {
				tmpl.ExecuteTemplate(w, "login.html", map[string]interface{}{"Error": "Invalid credentials"})
				return
			}
			http.SetCookie(w, &http.Cookie{
				Name:   "user_id",
				Value:  fmt.Sprint(userID),
				Path:   "/",
				MaxAge: 86400 * 30,
			})
			http.Redirect(w, r, "/profile", http.StatusSeeOther)
			return
		}
		tmpl.ExecuteTemplate(w, "login.html", nil)
	}
}
