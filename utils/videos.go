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
	rows, err := db.Query(`SELECT videos.id, videos.filename, videos.created_at, users.username FROM videos JOIN users ON videos.sender_id = users.id ORDER BY videos.created_at ASC`)
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
	_, err = db.Exec("DELETE FROM videos WHERE id = ?", id)
	if err != nil {
		return err
	}
	if filename != "" {
		_ = os.Remove("uploads/videos/" + filename)
	}
	return nil
}
