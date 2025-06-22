package utils

import (
	"database/sql"
	"os"
)

type Video struct {
	ID        int
	Filename  string
	CreatedAt string
	Sender    string
}

func AddVideo(db *sql.DB, senderID, filename string) error {
	_, err := db.Exec("INSERT INTO videos(sender_id, filename) VALUES (?, ?)", senderID, filename)
	return err
}

func GetAllVideos(db *sql.DB) ([]Video, error) {
	rows, err := db.Query(`SELECT v.id, v.filename, v.created_at, u.username FROM videos v JOIN users u ON v.sender_id = u.id ORDER BY v.id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var videos []Video
	for rows.Next() {
		var v Video
		if err := rows.Scan(&v.ID, &v.Filename, &v.CreatedAt, &v.Sender); err == nil {
			videos = append(videos, v)
		}
	}
	return videos, nil
}

func DeleteVideo(db *sql.DB, id string) error {
	var filename string
	err := db.QueryRow("SELECT filename FROM videos WHERE id=?", id).Scan(&filename)
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM videos WHERE id=?", id)
	if err != nil {
		return err
	}
	if filename != "" {
		_ = os.Remove("uploads/videos/" + filename)
	}
	return nil
}
