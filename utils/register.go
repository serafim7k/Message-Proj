package utils

import (
	"database/sql"
	"html/template"
	"net/http"
)

func RegisterHandler(db *sql.DB, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if r.Method == "POST" {
			username := r.FormValue("username")
			password := r.FormValue("password")
			_, err := db.Exec("INSERT INTO users(username, password) VALUES (?, ?)", username, password)
			if err != nil {
				tmpl.ExecuteTemplate(w, "register.html", map[string]interface{}{"Error": "Username taken"})
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
