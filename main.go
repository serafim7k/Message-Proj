package main

import (
	"GoWebSite/utils"
	"database/sql"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "root:1root2@tcp(127.0.0.1:3306)/myapp")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tmpl, err := template.ParseFiles(
		"templates/login.html",
		"templates/register.html",
		"templates/profile.html",
		"templates/index.html",
		"templates/header.html",
		"templates/footer.html",
		"templates/about.html",
		"templates/messenger.html",
		"templates/admin.html",
		"templates/photos.html",
		"templates/videos.html",
		"templates/music.html",
		"templates/pdfchat.html",
	)
	if err != nil {
		log.Fatal("Template parse error:", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
        id INT AUTO_INCREMENT PRIMARY KEY,
        username VARCHAR(255) NOT NULL UNIQUE,
        password VARCHAR(255) NOT NULL
    )`)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS messages (
		id INT AUTO_INCREMENT PRIMARY KEY,
		sender_id INT NOT NULL,
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (sender_id) REFERENCES users(id)
	)`)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS photos (
		id INT AUTO_INCREMENT PRIMARY KEY,
		sender_id INT NOT NULL,
		filename VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (sender_id) REFERENCES users(id)
	)`)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS videos (
		id INT AUTO_INCREMENT PRIMARY KEY,
		sender_id INT NOT NULL,
		filename VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (sender_id) REFERENCES users(id)
	)`)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS music (
		id INT AUTO_INCREMENT PRIMARY KEY,
		sender_id INT NOT NULL,
		filename VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (sender_id) REFERENCES users(id)
	)`)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS pdfs (
		id INT AUTO_INCREMENT PRIMARY KEY,
		sender_id INT NOT NULL,
		filename VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (sender_id) REFERENCES users(id)
	)`)

	// Добавить пользователя admin, если его нет
	_, err = db.Exec(`INSERT IGNORE INTO users(username, password) VALUES ('admin', 'dsa3hgkkh2138')`)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/uploads/photos/", http.StripPrefix("/uploads/photos/", http.FileServer(http.Dir("uploads/photos"))))
	http.Handle("/uploads/videos/", http.StripPrefix("/uploads/videos/", http.FileServer(http.Dir("uploads/videos"))))
	http.Handle("/uploads/music/", http.StripPrefix("/uploads/music/", http.FileServer(http.Dir("uploads/music"))))
	http.Handle("/uploads/pdfs/", http.StripPrefix("/uploads/pdfs/", http.FileServer(http.Dir("uploads/pdfs"))))

	http.HandleFunc("/register-form", func(w http.ResponseWriter, r *http.Request) {
		if cookie, err := r.Cookie("user_id"); err == nil && cookie.Value != "" {
			http.Redirect(w, r, "/profile", http.StatusSeeOther)
			return
		}
		utils.RegisterHandler(db, tmpl)(w, r)
	})

	http.HandleFunc("/login-form", func(w http.ResponseWriter, r *http.Request) {
		if cookie, err := r.Cookie("user_id"); err == nil && cookie.Value != "" {
			http.Redirect(w, r, "/profile", http.StatusSeeOther)
			return
		}
		utils.LoginHandler(db, tmpl)(w, r)
	})

	http.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		profileHandler(w, r, db, tmpl)
	})

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		if cookie, err := r.Cookie("user_id"); err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login-form", http.StatusSeeOther)
			return
		}
		logoutHandler(w, r)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		homeHandler(w, r, tmpl)
	})

	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("user_id")
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login-form", http.StatusSeeOther)
			return
		}
		isAdmin := utils.IsAdmin(db, cookie.Value)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.ExecuteTemplate(w, "about.html", map[string]interface{}{"LoggedIn": true, "IsAdmin": isAdmin})
	})

	http.HandleFunc("/messenger", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("user_id")
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login-form", http.StatusSeeOther)
			return
		}
		userID := cookie.Value
		var username string
		db.QueryRow("SELECT username FROM users WHERE id=?", userID).Scan(&username)
		isAdmin := utils.IsAdmin(db, userID)
		if r.Method == "POST" {
			content := r.FormValue("content")
			if content != "" {
				err := utils.AddMessage(db, userID, content)
				if err != nil {
					http.Error(w, "Failed to send message", http.StatusInternalServerError)
					return
				}
			}
			http.Redirect(w, r, "/messenger", http.StatusSeeOther)
			return
		}
		messages, err := utils.GetAllMessages(db)
		if err != nil {
			http.Error(w, "Failed to load messages", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.ExecuteTemplate(w, "messenger.html", map[string]interface{}{"Messages": messages, "LoggedIn": true, "IsAdmin": isAdmin})
	})

	http.HandleFunc("/messenger/delete", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("user_id")
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login-form", http.StatusSeeOther)
			return
		}
		var username string
		err = db.QueryRow("SELECT username FROM users WHERE id=?", cookie.Value).Scan(&username)
		if err != nil || username != "admin" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		if r.Method == "POST" {
			msgID := r.FormValue("id")
			if msgID != "" {
				_ = utils.DeleteMessage(db, msgID)
			}
		}
		http.Redirect(w, r, "/messenger", http.StatusSeeOther)
	})

	http.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("user_id")
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login-form", http.StatusSeeOther)
			return
		}
		var username string
		err = db.QueryRow("SELECT username FROM users WHERE id=?", cookie.Value).Scan(&username)
		if err != nil || username != "admin" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		users, err := utils.GetAllUsers(db)
		if err != nil {
			http.Error(w, "Failed to load users", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.ExecuteTemplate(w, "admin.html", map[string]interface{}{"Users": users, "LoggedIn": true, "IsAdmin": true})
	})

	http.HandleFunc("/admin/delete", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("user_id")
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login-form", http.StatusSeeOther)
			return
		}
		var username string
		err = db.QueryRow("SELECT username FROM users WHERE id=?", cookie.Value).Scan(&username)
		if err != nil || username != "admin" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		if r.Method == "POST" {
			userID := r.FormValue("id")
			if userID != "" {
				_ = utils.DeleteUser(db, userID)
			}
		}
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	})

	http.HandleFunc("/photos", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("user_id")
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login-form", http.StatusSeeOther)
			return
		}
		userID := cookie.Value
		isAdmin := utils.IsAdmin(db, userID)
		if r.Method == "POST" {
			r.ParseMultipartForm(10 << 20)
			file, handler, err := r.FormFile("photo")
			if err == nil && handler != nil {
				defer file.Close()
				os.MkdirAll("uploads/photos", 0755)
				f, err := os.Create("uploads/photos/" + handler.Filename)
				if err == nil {
					defer f.Close()
					io.Copy(f, file)
					utils.AddPhoto(db, userID, handler.Filename)
				}
			}
			http.Redirect(w, r, "/photos", http.StatusSeeOther)
			return
		}
		photos, err := utils.GetAllPhotos(db)
		if err != nil {
			http.Error(w, "Failed to load photos", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.ExecuteTemplate(w, "photos.html", map[string]interface{}{"Photos": photos, "LoggedIn": true, "IsAdmin": isAdmin})
	})

	http.HandleFunc("/photos/delete", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/photos/delete called")
		cookie, err := r.Cookie("user_id")
		if err != nil || cookie.Value == "" || !utils.IsAdmin(db, cookie.Value) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		if r.Method == "POST" {
			photoID := r.FormValue("id")
			log.Println("photoID:", photoID)
			if photoID != "" {
				_ = utils.DeletePhoto(db, photoID)
			}
		}
		http.Redirect(w, r, "/photos", http.StatusSeeOther)
	})

	http.HandleFunc("/videos", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("user_id")
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login-form", http.StatusSeeOther)
			return
		}
		userID := cookie.Value
		isAdmin := utils.IsAdmin(db, userID)
		if r.Method == "POST" {
			r.ParseMultipartForm(100 << 20)
			file, handler, err := r.FormFile("video")
			if err == nil && handler != nil {
				defer file.Close()
				os.MkdirAll("uploads/videos", 0755)
				f, err := os.Create("uploads/videos/" + handler.Filename)
				if err == nil {
					defer f.Close()
					io.Copy(f, file)
					utils.AddVideo(db, userID, handler.Filename)
				}
			}
			http.Redirect(w, r, "/videos", http.StatusSeeOther)
			return
		}
		videos, err := utils.GetAllVideos(db)
		if err != nil {
			http.Error(w, "Failed to load videos", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.ExecuteTemplate(w, "videos.html", map[string]interface{}{"Videos": videos, "LoggedIn": true, "IsAdmin": isAdmin})
	})

	http.HandleFunc("/videos/delete", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/videos/delete called")
		cookie, err := r.Cookie("user_id")
		if err != nil || cookie.Value == "" || !utils.IsAdmin(db, cookie.Value) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		if r.Method == "POST" {
			videoID := r.FormValue("id")
			log.Println("videoID:", videoID)
			if videoID != "" {
				_ = utils.DeleteVideo(db, videoID)
			}
		}
		http.Redirect(w, r, "/videos", http.StatusSeeOther)
	})

	http.HandleFunc("/music", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("user_id")
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login-form", http.StatusSeeOther)
			return
		}
		userID := cookie.Value
		isAdmin := utils.IsAdmin(db, userID)
		if r.Method == "POST" {
			r.ParseMultipartForm(20 << 20)
			file, handler, err := r.FormFile("music")
			if err == nil && handler != nil {
				defer file.Close()
				os.MkdirAll("uploads/music", 0755)
				f, err := os.Create("uploads/music/" + handler.Filename)
				if err == nil {
					defer f.Close()
					io.Copy(f, file)
					utils.AddMusic(db, userID, handler.Filename)
				}
			}
			http.Redirect(w, r, "/music", http.StatusSeeOther)
			return
		}
		music, err := utils.GetAllMusic(db)
		if err != nil {
			http.Error(w, "Failed to load music", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.ExecuteTemplate(w, "music.html", map[string]interface{}{"Music": music, "LoggedIn": true, "IsAdmin": isAdmin})
	})

	http.HandleFunc("/music/delete", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("user_id")
		if err != nil || cookie.Value == "" || !utils.IsAdmin(db, cookie.Value) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		if r.Method == "POST" {
			musicID := r.FormValue("id")
			if musicID != "" {
				_ = utils.DeleteMusic(db, musicID)
			}
		}
		http.Redirect(w, r, "/music", http.StatusSeeOther)
	})

	http.HandleFunc("/pdfchat", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("user_id")
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login-form", http.StatusSeeOther)
			return
		}
		userID := cookie.Value
		isAdmin := utils.IsAdmin(db, userID)
		if r.Method == "POST" {
			r.ParseMultipartForm(10 << 20)
			file, handler, err := r.FormFile("pdf")
			if err == nil && handler != nil {
				defer file.Close()
				os.MkdirAll("uploads/pdfs", 0755)
				f, err := os.Create("uploads/pdfs/" + handler.Filename)
				if err == nil {
					defer f.Close()
					io.Copy(f, file)
					utils.AddPDF(db, userID, handler.Filename)
				}
			}
			http.Redirect(w, r, "/pdfchat", http.StatusSeeOther)
			return
		}
		pdfs, err := utils.GetAllPDFs(db)
		if err != nil {
			http.Error(w, "Failed to load pdfs", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.ExecuteTemplate(w, "pdfchat.html", map[string]interface{}{"PDFs": pdfs, "LoggedIn": true, "IsAdmin": isAdmin})
	})

	http.HandleFunc("/pdfchat/delete", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("user_id")
		if err != nil || cookie.Value == "" || !utils.IsAdmin(db, cookie.Value) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		if r.Method == "POST" {
			pdfID := r.FormValue("id")
			if pdfID != "" {
				_ = utils.DeletePDF(db, pdfID)
			}
		}
		http.Redirect(w, r, "/pdfchat", http.StatusSeeOther)
	})

	log.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	cookie, err := r.Cookie("user_id")
	loggedIn := err == nil
	isAdmin := false
	if loggedIn {
		isAdmin = utils.IsAdmin(db, cookie.Value)
	}
	if err := tmpl.ExecuteTemplate(w, "index.html", map[string]interface{}{"LoggedIn": loggedIn, "IsAdmin": isAdmin}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func profileHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, tmpl *template.Template) {
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Redirect(w, r, "/login-form", http.StatusSeeOther)
		return
	}
	var username string
	err = db.QueryRow("SELECT username FROM users WHERE id=?", cookie.Value).Scan(&username)
	isAdmin := utils.IsAdmin(db, cookie.Value)
	if err != nil {
		http.Error(w, "Пользователь не найден", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "profile.html", map[string]interface{}{"Username": username, "LoggedIn": true, "IsAdmin": isAdmin}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "user_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
