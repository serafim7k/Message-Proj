package utils

import (
	"database/sql"
	"html/template"
	"net/http"
	"strings"
)

func RegisterHandler(db *sql.DB, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if r.Method == "POST" {
			username := strings.ToLower(r.FormValue("username"))
			password := r.FormValue("password")
			var exists int
			err := db.QueryRow("SELECT COUNT(*) FROM users WHERE LOWER(username) = $1", username).Scan(&exists)
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
			if exists > 0 {
				tmpl.ExecuteTemplate(w, "register.html", map[string]interface{}{"Error": "Username taken"})
				return
			}
			_, err = db.Exec("INSERT INTO users(username, password) VALUES ($1, $2)", username, password)
			if err != nil {
				tmpl.ExecuteTemplate(w, "register.html", map[string]interface{}{"Error": "Registration error"})
				return
			}
			http.Redirect(w, r, "/login-form", http.StatusSeeOther)
			return
		}
		if err := tmpl.ExecuteTemplate(w, "register.html", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
