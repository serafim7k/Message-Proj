package utils

import (
	"database/sql"
	"os"
)

type Photo struct {
	ID        int
	Filename  string
	CreatedAt string
	Sender    string
}

func AddPhoto(db *sql.DB, senderID, filename string) error {
	_, err := db.Exec("INSERT INTO photos(sender_id, filename) VALUES (?, ?)", senderID, filename)
	return err
}

func GetAllPhotos(db *sql.DB) ([]Photo, error) {
	rows, err := db.Query(`SELECT photos.id, photos.filename, photos.created_at, users.username FROM photos JOIN users ON photos.sender_id = users.id ORDER BY photos.created_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var photos []Photo
	for rows.Next() {
		var p Photo
		if err := rows.Scan(&p.ID, &p.Filename, &p.CreatedAt, &p.Sender); err == nil {
			photos = append(photos, p)
		}
	}
	return photos, nil
}

func DeletePhoto(db *sql.DB, id string) error {
	var filename string
	err := db.QueryRow("SELECT filename FROM photos WHERE id=?", id).Scan(&filename)
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM photos WHERE id = ?", id)
	if err != nil {
		return err
	}
	if filename != "" {
		_ = os.Remove("uploads/photos/" + filename)
	}
	return nil
}
