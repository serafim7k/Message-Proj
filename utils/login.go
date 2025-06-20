package utils

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
)

func LoginHandler(db *sql.DB, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			username := r.FormValue("username")
			password := r.FormValue("password")
			var id int
			err := db.QueryRow("SELECT id FROM users WHERE username=? AND password=?", username, password).Scan(&id)
			if err != nil {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				tmpl.ExecuteTemplate(w, "login.html", map[string]interface{}{"Error": "Invalid credentials"})
				return
			}
			http.SetCookie(w, &http.Cookie{
				Name:  "user_id",
				Value: fmt.Sprint(id),
				Path:  "/",
			})
			http.Redirect(w, r, "/profile", http.StatusSeeOther)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.ExecuteTemplate(w, "login.html", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
