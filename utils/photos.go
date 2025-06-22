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
	_, err := db.Exec("INSERT INTO photos(sender_id, filename) VALUES ($1, $2)", senderID, filename)
	return err
}

func GetAllPhotos(db *sql.DB) ([]Photo, error) {
	rows, err := db.Query(`SELECT p.id, p.filename, p.created_at, u.username FROM photos p JOIN users u ON p.sender_id = u.id ORDER BY p.id ASC`)
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
	err := db.QueryRow("SELECT filename FROM photos WHERE id=$1", id).Scan(&filename)
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM photos WHERE id=$1", id)
	if err != nil {
		return err
	}
	if filename != "" {
		_ = os.Remove("uploads/photos/" + filename)
	}
	return nil
}
